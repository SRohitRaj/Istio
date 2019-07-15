package v2

import (
	"sync"

	"istio.io/istio/pilot/pkg/model"
)

type PushInformation struct {
	// If not empty, it is used to indicate the event is caused by a change in the clusters.
	// Only EDS for the listed clusters will be sent.
	edsUpdatedServices map[string]struct{}

	push *model.PushContext

	full bool
}

type PushQueue struct {
	mu          *sync.RWMutex
	cond        *sync.Cond
	connections map[*XdsConnection]*PushInformation
	order       []*XdsConnection
}

func NewPushQueue() *PushQueue {
	mu := &sync.RWMutex{}
	return &PushQueue{
		mu:          mu,
		connections: make(map[*XdsConnection]*PushInformation),
		cond:        sync.NewCond(mu),
	}
}

// Add will mark a proxy as pending a push. If it is already pending, pushInfo will be merged.
// edsUpdatedServices will be added together, and full will be set if either were full
func (p *PushQueue) Enqueue(proxy *XdsConnection, pushInfo *PushInformation) {

	p.mu.Lock()
	defer p.mu.Unlock()
	info, exists := p.connections[proxy]
	if !exists {
		p.connections[proxy] = pushInfo
		p.order = append(p.order, proxy)
	} else {
		info.push = pushInfo.push
		info.full = info.full || pushInfo.full

		edsUpdates := map[string]struct{}{}
		for endpoint := range pushInfo.edsUpdatedServices {
			edsUpdates[endpoint] = struct{}{}
		}
		for endpoint := range info.edsUpdatedServices {
			edsUpdates[endpoint] = struct{}{}
		}
		info.edsUpdatedServices = edsUpdates
	}
	p.cond.Broadcast()
}

// Remove a proxy from the queue. If there are no proxies ready to be removed, this will block
func (p *PushQueue) Dequeue() (*XdsConnection, *PushInformation) {
	p.mu.Lock()
	if len(p.order) == 0 {
		p.cond.Wait()
	}

	defer p.mu.Unlock()
	head := p.order[0]
	p.order = p.order[1:]
	info := p.connections[head]
	delete(p.connections, head)
	return head, info
}

// Get number of pending proxies
func (p *PushQueue) Pending() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.order)
}
