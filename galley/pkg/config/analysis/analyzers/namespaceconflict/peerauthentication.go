package namespaceconflict

import (
	"fmt"


	k8s_labels "k8s.io/apimachinery/pkg/labels"

	v1beta1 "istio.io/api/security/v1beta1"
	"istio.io/istio/galley/pkg/config/analysis"
	"istio.io/istio/galley/pkg/config/analysis/msg"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

// Analyzer checks conditions related to conflicting namespace level resources
type PeerAuthenticationConflictAnalyzer struct{}

var _ analysis.Analyzer = &PeerAuthenticationConflictAnalyzer{}

var peerAuthenticationCol = collections.IstioSecurityV1Beta1Peerauthentications.Name()

// Metadata implements Analyzer
func (a *PeerAuthenticationConflictAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "namespaceconflict.PeerAuthenticationConflictAnalyzer",
		Description: "Checks conditions related to Peer Authentication conflicting namespace level resources",
		Inputs: collection.Names{
			peerAuthenticationCol,
		},
	}
}

func (a *PeerAuthenticationConflictAnalyzer) Analyze(c analysis.Context) {
	namespaceWideConfiguration := map[string]*resource.Instance{}

	c.ForEach(peerAuthenticationCol, func(r *resource.Instance) bool {
		x := r.Message.(*v1beta1.PeerAuthentication)
		xNS := r.Metadata.FullName.Namespace.String()

		// If the resource has workloads associated with it, analyze for conflicts with selector
		if x.GetSelector() != nil {
			// If there's namespace wide configuration, there's a conflict.
			if val, found := namespaceWideConfiguration[xNS]; found {
				conflicts := []string{val.Metadata.FullName.String(), r.Metadata.FullName.String()}
				m := msg.NewNamespaceResourceConflict(r, peerAuthenticationCol.String(), xNS, fmt.Sprintf("(ALL) Namespace: %v", xNS), conflicts)
				c.Report(collections.IstioSecurityV1Beta1Peerauthentications.Name(), m)
			} else {
				a.analyzeWorkloadSelectorConflicts(r, c)
			}
		} else {
			// If there's namespace wide configuration, there's a conflict.
			if val, found := namespaceWideConfiguration[xNS]; found {
				conflicts := []string{val.Metadata.FullName.String(), r.Metadata.FullName.String()}
				m := msg.NewNamespaceResourceConflict(r, peerAuthenticationCol.String(), xNS, fmt.Sprintf("(ALL) Namespace: %v", xNS), conflicts)
				c.Report(collections.IstioSecurityV1Beta1Peerauthentications.Name(), m)
			} else {
				namespaceWideConfiguration[xNS] = r
			}
		}
		return true
	})
}

func (a *PeerAuthenticationConflictAnalyzer) analyzeWorkloadSelectorConflicts(r *resource.Instance, c analysis.Context) {
	x := r.Message.(*v1beta1.PeerAuthentication)
	xNS := r.Metadata.FullName.Namespace.String()

	// Find all resources that have the same selector
	matches := a.findMatchingSelectors(r, c)

	// There should be only one resource associated with a selector
	if len(matches) != 0 {
		conflicts := []string{}
		for _, match := range matches {
			conflicts = append(conflicts, match.Metadata.FullName.String())
		}
		m := msg.NewNamespaceResourceConflict(r, peerAuthenticationCol.String(), xNS, k8s_labels.SelectorFromSet(x.GetSelector().MatchLabels).String(), conflicts)
		c.Report(collections.IstioSecurityV1Beta1Peerauthentications.Name(), m)
		return
	}
}

// Finds all resources that have the same selector as the resource we're checking
func (a *PeerAuthenticationConflictAnalyzer) findMatchingSelectors(r *resource.Instance, c analysis.Context) []*resource.Instance {
	x := r.Message.(*v1beta1.PeerAuthentication)
	xName := r.Metadata.FullName.String()
	xNS := r.Metadata.FullName.Namespace.String()
	xSelector := k8s_labels.SelectorFromSet(x.GetSelector().MatchLabels).String()
	fmt.Println(xSelector)
	matches := []*resource.Instance{}
	c.ForEach(peerAuthenticationCol, func(r1 *resource.Instance) bool {
		y := r1.Message.(*v1beta1.PeerAuthentication)
		yName := r1.Metadata.FullName.String()
		yNS := r1.Metadata.FullName.Namespace.String()
		ySelector := k8s_labels.SelectorFromSet(y.GetSelector().MatchLabels).String()
		fmt.Println(ySelector)
		if xSelector == ySelector && xName != yName && xNS == yNS {
			matches = append(matches, r)
			matches = append(matches, r1)
		}

		return true
	})
	return matches
}
