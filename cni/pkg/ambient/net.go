// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ambient

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/netip"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"golang.org/x/sys/unix"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	"istio.io/istio/cni/pkg/ambient/constants"
	ebpf "istio.io/istio/cni/pkg/ebpf/server"
	pconstants "istio.io/istio/pkg/config/constants"
	istiolog "istio.io/pkg/log"
)

var log = istiolog.RegisterScope("ambient", "ambient controller", 0)

func IsPodInIpset(pod *corev1.Pod) bool {
	ipset, err := Ipset.List()
	if err != nil {
		log.Errorf("Failed to list ipset entries: %v", err)
		return false
	}

	// Since not all kernels support comments in ipset, we should also try and
	// match against the IP
	for _, ip := range ipset {
		if ip.Comment == string(pod.UID) {
			return true
		}
		if ip.IP.String() == pod.Status.PodIP {
			return true
		}
	}

	return false
}

func RouteExists(rte []string) bool {
	output, err := executeOutput(
		"bash", "-c",
		fmt.Sprintf("ip route show %s | wc -l", strings.Join(rte, " ")),
	)
	if err != nil {
		return false
	}

	log.Debugf("RouteExists(%s): %s", strings.Join(rte, " "), output)

	return output == "1"
}

func AddPodToMesh(client kubernetes.Interface, pod *corev1.Pod, ip string) {
	addPodToMeshWithIptables(client, pod, ip)
}

func addPodToMeshWithIptables(client kubernetes.Interface, pod *corev1.Pod, ip string) {
	if ip == "" {
		ip = pod.Status.PodIP
	}
	if ip == "" {
		log.Debugf("skip adding pod %s/%s, IP not yet allocated", pod.Name, pod.Namespace)
		return
	}

	if !IsPodInIpset(pod) {
		log.Infof("Adding pod '%s/%s' (%s) to ipset", pod.Name, pod.Namespace, string(pod.UID))
		err := Ipset.AddIP(net.ParseIP(ip).To4(), string(pod.UID))
		if err != nil {
			log.Errorf("Failed to add pod %s to ipset list: %v", pod.Name, err)
		}
	} else {
		log.Infof("Pod '%s/%s' (%s) is in ipset", pod.Name, pod.Namespace, string(pod.UID))
	}

	rte, err := buildRouteFromPod(pod, ip)
	if err != nil {
		log.Errorf("Failed to build route for pod %s: %v", pod.Name, err)
	}

	if !RouteExists(rte) {
		log.Infof("Adding route for %s/%s: %+v", pod.Name, pod.Namespace, rte)
		// @TODO Try and figure out why buildRouteFromPod doesn't return a good route that we can
		// use err = netlink.RouteAdd(rte):
		// Error: {"level":"error","time":"2022-06-24T16:30:59.083809Z","msg":"Failed to add route ({Ifindex: 4 Dst: 10.244.2.7/32
		// Via: Family: 2, Address: 192.168.126.2 Src: 10.244.2.1 Gw: <nil> Flags: [] Table: 100 Realm: 0}) for pod
		// helloworld-v2-same-node-67b6b764bf-zhmp4: invalid argument"}
		err = execute("ip", append([]string{"route", "add"}, rte...)...)
		if err != nil {
			log.Warnf("Failed to add route (%s) for pod %s: %v", rte, pod.Name, err)
		}
	} else {
		log.Infof("Route already exists for %s/%s: %+v", pod.Name, pod.Namespace, rte)
	}

	dev, err := getDeviceWithDestinationOf(ip)
	if err != nil {
		log.Warnf("Failed to get device for destination %s", ip)
		return
	}
	err = SetProc("/proc/sys/net/ipv4/conf/"+dev+"/rp_filter", "0")
	if err != nil {
		log.Warnf("Failed to set rp_filter to 0 for device %s", dev)
	}

	if err := AnnotateEnrolledPod(client, pod); err != nil {
		log.Errorf("failed to annotate pod enrollment: %v", err)
	}
}

var annotationPatch = []byte(fmt.Sprintf(
	`{"metadata":{"annotations":{"%s":"%s"}}}`,
	pconstants.AmbientRedirection,
	pconstants.AmbientRedirectionEnabled,
))

var annotationRemovePatch = []byte(fmt.Sprintf(
	`{"metadata":{"annotations":{"%s":null}}}`,
	pconstants.AmbientRedirection,
))

func AnnotateEnrolledPod(client kubernetes.Interface, pod *corev1.Pod) error {
	_, err := client.CoreV1().
		Pods(pod.Namespace).
		Patch(
			context.Background(),
			pod.Name,
			types.MergePatchType,
			annotationPatch,
			metav1.PatchOptions{},
		)
	return err
}

func AnnotateUnenrollPod(client kubernetes.Interface, pod *corev1.Pod) error {
	if pod.Annotations[pconstants.AmbientRedirection] != pconstants.AmbientRedirectionEnabled {
		return nil
	}
	// TODO: do not overwrite if already none
	_, err := client.CoreV1().
		Pods(pod.Namespace).
		Patch(
			context.Background(),
			pod.Name,
			types.MergePatchType,
			annotationRemovePatch,
			metav1.PatchOptions{},
		)
	return err
}

func DelPodFromMesh(client kubernetes.Interface, pod *corev1.Pod) {
	log.Debugf("Removing pod '%s/%s' (%s) from mesh", pod.Name, pod.Namespace, string(pod.UID))
	if IsPodInIpset(pod) {
		log.Infof("Removing pod '%s' (%s) from ipset", pod.Name, string(pod.UID))
		err := Ipset.DeleteIP(net.ParseIP(pod.Status.PodIP).To4())
		if err != nil {
			log.Errorf("Failed to delete pod %s from ipset list: %v", pod.Name, err)
		}
	} else {
		log.Infof("Pod '%s/%s' (%s) is not in ipset", pod.Name, pod.Namespace, string(pod.UID))
	}
	rte, err := buildRouteFromPod(pod, "")
	if err != nil {
		log.Errorf("Failed to build route for pod %s: %v", pod.Name, err)
	}
	if RouteExists(rte) {
		log.Infof("Removing route: %+v", rte)
		// @TODO Try and figure out why buildRouteFromPod doesn't return a good route that we can
		// use this:
		// err = netlink.RouteDel(rte)
		err = execute("ip", append([]string{"route", "del"}, rte...)...)
		if err != nil {
			log.Warnf("Failed to delete route (%s) for pod %s: %v", rte, pod.Name, err)
		}
	}

	if err := AnnotateUnenrollPod(client, pod); err != nil {
		log.Errorf("failed to annotate pod unenrollment: %v", err)
	}
}

func buildRouteFromPod(pod *corev1.Pod, ip string) ([]string, error) {
	if ip == "" {
		ip = pod.Status.PodIP
	}

	if ip == "" {
		return nil, errors.New("no ip found")
	}

	return []string{
		"table",
		fmt.Sprintf("%d", constants.RouteTableInbound),
		fmt.Sprintf("%s/32", ip),
		"via",
		constants.ZTunnelInboundTunIP,
		"dev",
		constants.InboundTun,
		"src",
		HostIP,
	}, nil
}

func buildEbpfArgsByIP(ip string, isZtuunel, isRemove bool) (*ebpf.RedirectArgs, error) {
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ip(%s): %v", ip, err)
	}
	veth, err := getVethWithDestinationOf(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}
	peerIndex, err := netlink.VethPeerIndex(veth)
	if err != nil {
		return nil, fmt.Errorf("failed to get veth peerIndex: %v", err)
	}

	peerNs, err := getNsNameFromNsID(veth.Attrs().NetNsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ns name: %v", err)
	}

	mac, err := getMacFromNsIdx(peerNs, peerIndex)
	if err != nil {
		return nil, err
	}

	return &ebpf.RedirectArgs{
		IPAddrs:   []netip.Addr{ipAddr},
		MacAddr:   mac,
		Ifindex:   veth.Attrs().Index,
		PeerIndex: peerIndex,
		PeerNs:    peerNs,
		IsZtunnel: isZtuunel,
		Remove:    isRemove,
	}, nil
}

func GetIndexAndPeerMac(podIfName, ns string) (int, net.HardwareAddr, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	curNs, err := netns.Get()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get cur nshandler: %v", err)
	}
	defer func() {
		if err := curNs.Close(); err != nil {
			log.Errorf("close ns handler failure: %v", err)
		}
	}()

	ns = filepath.Base(ns)
	nsHdlr, err := netns.GetFromName(ns)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get ns(%s) handler: %v", ns, err)
	}
	defer nsHdlr.Close()

	err = netns.Set(nsHdlr)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to switch to net ns(%s): %v", ns, err)
	}
	defer func() {
		if err := netns.Set(curNs); err != nil {
			log.Errorf("set back ns failure: %v", err)
		}
	}()

	link, err := netlink.LinkByName(podIfName)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get link(%s) in ns(%s): %v", podIfName, ns, err)
	}

	veth, ok := link.(*netlink.Veth)
	if !ok {
		return 0, nil, errors.New("not veth implemented CNI")
	}

	hostIfIndex, err := netlink.VethPeerIndex(veth)
	if err != nil {
		return 0, nil, err
	}

	return hostIfIndex, veth.Attrs().HardwareAddr, nil
}

func getMacFromNsIdx(ns string, ifIndex int) (net.HardwareAddr, error) {
	nsHdlr, err := netns.GetFromName(ns)
	if err != nil {
		return nil, fmt.Errorf("failed to get ns(%s) handler: %v", ns, err)
	}
	defer nsHdlr.Close()
	nl, err := netlink.NewHandleAt(nsHdlr)
	if err != nil {
		return nil, fmt.Errorf("failed to link handler for ns(%s): %v", ns, err)
	}
	defer nl.Close()
	link, err := nl.LinkByIndex(ifIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get link(%d) in ns(%s): %v", ifIndex, ns, err)
	}
	return link.Attrs().HardwareAddr, nil
}

func getNsNameFromNsID(nsid int) (string, error) {
	foundNs := errors.New("nsid found, stop iterating")
	nsName := ""
	err := filepath.WalkDir("/var/run/netns", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		fd, err := unix.Open(p, unix.O_RDONLY, 0)
		if err != nil {
			log.Warnf("failed to open: %v", err)
			return nil
		}
		defer unix.Close(fd)

		id, err := netlink.GetNetNsIdByFd(fd)
		if err != nil {
			log.Warnf("failed to open: %v", err)
			return nil
		}
		if id == nsid {
			nsName = path.Base(p)
			return foundNs
		}
		return nil
	})
	if err == foundNs {
		return nsName, nil
	}
	return "", fmt.Errorf("failed to get namespace for %d", nsid)
}

func getLinkWithDestinationOf(ip string) (netlink.Link, error) {
	routes, err := netlink.RouteListFiltered(
		netlink.FAMILY_V4,
		&netlink.Route{Dst: &net.IPNet{IP: net.ParseIP(ip), Mask: net.CIDRMask(32, 32)}},
		netlink.RT_FILTER_DST)
	if err != nil {
		return nil, err
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no routes found for %s", ip)
	}

	linkIndex := routes[0].LinkIndex
	return netlink.LinkByIndex(linkIndex)
}

func getVethWithDestinationOf(ip string) (*netlink.Veth, error) {
	link, err := getLinkWithDestinationOf(ip)
	if err != nil {
		return nil, err
	}
	veth, ok := link.(*netlink.Veth)
	if !ok {
		return nil, errors.New("not veth implemented CNI")
	}
	return veth, nil
}

func getDeviceWithDestinationOf(ip string) (string, error) {
	link, err := getLinkWithDestinationOf(ip)
	if err != nil {
		return "", err
	}
	return link.Attrs().Name, nil
}

// GetHostIPByRoute get the automatically chosen host ip to the Pod's CIDR
func GetHostIPByRoute(kubeClient kubernetes.Interface) (string, error) {
	// We assume per node POD's CIDR is the same block, so the route to the POD
	// from host should be "same". Otherwise, there may multiple host IPs will be
	// used as source to dial to PODs.
	pods, err := kubeClient.CoreV1().Pods(metav1.NamespaceAll).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: "app=ztunnel",
			FieldSelector: "spec.nodeName=" + NodeName,
		})
	if err != nil {
		return "", fmt.Errorf("error getting ztunnel node: %v", err)
	}
	for _, pod := range pods.Items {
		targetIP := pod.Status.PodIP
		if hostIP := getOutboundIP(targetIP); hostIP != nil {
			return hostIP.String(), nil
		}
	}
	return "", fmt.Errorf("failed to get outbound IP to Pods")
}

// Get preferred outbound ip of this machine
func getOutboundIP(ip string) net.IP {
	conn, err := net.Dial("udp", ip+":80")
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GetHostIP(kubeClient kubernetes.Interface) (string, error) {
	var ip string
	// Get the node from the Kubernetes API
	node, err := kubeClient.CoreV1().Nodes().Get(context.TODO(), NodeName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error getting node: %v", err)
	}

	ip = node.Spec.PodCIDR

	// This needs to be done as in Kind, the node internal IP is not the one we want.
	if ip == "" {
		// PodCIDR is not set, try to get the IP from the node internal IP
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				return address.Address, nil
			}
		}
	} else {
		network, err := netip.ParsePrefix(ip)
		if err != nil {
			return "", fmt.Errorf("error parsing node IP: %v", err)
		}

		ifaces, err := net.Interfaces()
		if err != nil {
			return "", fmt.Errorf("error getting interfaces: %v", err)
		}

		for _, iface := range ifaces {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", fmt.Errorf("error getting addresses: %v", err)
			}

			for _, addr := range addrs {
				a, err := netip.ParseAddr(strings.Split(addr.String(), "/")[0])
				if err != nil {
					return "", fmt.Errorf("error parsing address: %v", err)
				}
				if network.Contains(a) {
					return a.String(), nil
				}
			}
		}
	}
	// fall back to use Node Internal IP
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			return address.Address, nil
		}
	}
	return "", nil
}

// CreateRulesOnNode initializes the routing, firewall and ipset rules on the node.
func (s *Server) CreateRulesOnNode(ztunnelVeth, ztunnelIP string, captureDNS bool) error {
	var err error

	log.Debugf("CreateRulesOnNode: ztunnelVeth=%s, ztunnelIP=%s", ztunnelVeth, ztunnelIP)

	// Check if chain exists, if it exists flush.. otherwise initialize
	err = execute(s.IptablesCmd(), "-t", "mangle", "-C", "output", "-j", constants.ChainZTunnelOutput)
	if err == nil {
		log.Debugf("Chain %s already exists, flushing", constants.ChainOutput)
		s.flushLists()
	} else {
		log.Debugf("Initializing lists")
		err = s.initializeLists()
		if err != nil {
			return err
		}
	}

	// Create ipset of pod members.
	log.Debug("Creating ipset")
	err = Ipset.CreateSet()
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("error creating ipset: %v", err)
	}

	appendRules := []*iptablesRule{
		// Skip things that come from the tunnels, but don't apply the conn skip mark
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-i", constants.InboundTun,
			"-j", "MARK",
			"--set-mark", constants.SkipMark,
		),
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-i", constants.InboundTun,
			"-j", "RETURN",
		),
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-i", constants.OutboundTun,
			"-j", "MARK",
			"--set-mark", constants.SkipMark,
		),
		newIptableRule(constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-i", constants.OutboundTun,
			"-j", "RETURN",
		),

		// Make sure that whatever is skipped is also skipped for returning packets.
		// If we have a skip mark, save it to conn mark.
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelForward,
			"-m", "mark",
			"--mark", constants.ConnSkipMark,
			"-j", "CONNMARK",
			"--save-mark",
			"--nfmask", constants.ConnSkipMask,
			"--ctmask", constants.ConnSkipMask,
		),
		// Input chain might be needed for things in host namespace that are skipped.
		// Place the mark here after routing was done, not sure if conn-tracking will figure
		// it out if I do it before, as NAT might change the connection tuple.
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelInput,
			"-m", "mark",
			"--mark", constants.ConnSkipMark,
			"-j", "CONNMARK",
			"--save-mark",
			"--nfmask", constants.ConnSkipMask,
			"--ctmask", constants.ConnSkipMask,
		),

		// For things with the proxy mark, we need different routing just on returning packets
		// so we give a different mark to them.
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelForward,
			"-m", "mark",
			"--mark", constants.ProxyMark,
			"-j", "CONNMARK",
			"--save-mark",
			"--nfmask", constants.ProxyMask,
			"--ctmask", constants.ProxyMask,
		),
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelInput,
			"-m", "mark",
			"--mark", constants.ProxyMark,
			"-j", "CONNMARK",
			"--save-mark",
			"--nfmask", constants.ProxyMask,
			"--ctmask", constants.ProxyMask,
		),
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelOutput,
			"--source", HostIP,
			"-j", "MARK",
			"--set-mark", constants.ConnSkipMask,
		),

		// If we have an outbound mark, we don't need kube-proxy to do anything,
		// so accept it before kube-proxy translates service vips to pod ips
		newIptableRule(
			constants.TableNat,
			constants.ChainZTunnelPrerouting,
			"-m", "mark",
			"--mark", constants.OutboundMark,
			"-j", "ACCEPT",
		),
		newIptableRule(
			constants.TableNat,
			constants.ChainZTunnelPostrouting,
			"-m", "mark",
			"--mark", constants.OutboundMark,
			"-j", "ACCEPT",
		),
	}

	if captureDNS {
		appendRules = append(appendRules,
			newIptableRule(
				constants.TableNat,
				constants.ChainZTunnelPrerouting,
				"-p", "udp",
				"-m", "set",
				"--match-set", Ipset.Name, "src",
				"--dport", "53",
				"-j", "DNAT",
				"--to", fmt.Sprintf("%s:%d", ztunnelIP, constants.DNSCapturePort),
			),
		)
	}

	appendRules2 := []*iptablesRule{
		// Don't set anything on the tunnel (geneve port is 6081), as the tunnel copies
		// the mark to the un-tunneled packet.
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-p", "udp",
			"-m", "udp",
			"--dport", "6081",
			"-j", "RETURN",
		),

		// If we have the conn mark, restore it to mark, to make sure that the other side of the connection
		// is skipped as well.
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-m", "connmark",
			"--mark", constants.ConnSkipMark,
			"-j", "MARK",
			"--set-mark", constants.SkipMark,
		),
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-m", "mark",
			"--mark", constants.SkipMark,
			"-j", "RETURN",
		),

		// If we have the proxy mark in, set the return mark to make sure that original src packets go to ztunnel
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"!", "-i", ztunnelVeth,
			"-m", "connmark",
			"--mark", constants.ProxyMark,
			"-j", "MARK",
			"--set-mark", constants.ProxyRetMark,
		),
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-m", "mark",
			"--mark", constants.ProxyRetMark,
			"-j", "RETURN",
		),

		// Send fake source outbound connections to the outbound route table (for original src)
		// if it's original src, the source ip of packets coming from the proxy might be that of a pod, so
		// make sure we don't tproxy it
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-i", ztunnelVeth,
			"!", "--source", ztunnelIP,
			"-j", "MARK",
			"--set-mark", constants.ProxyMark,
		),
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-m", "mark",
			"--mark", constants.SkipMark,
			"-j", "RETURN",
		),

		// Make sure anything that leaves ztunnel is routed normally (xds, connections to other ztunnels,
		// connections to upstream pods...)
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-i", ztunnelVeth,
			"-j", "MARK",
			"--set-mark", constants.ConnSkipMark,
		),

		// skip udp so DNS works. We can make this more granular.
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-p", "udp",
			"-j", "MARK",
			"--set-mark", constants.ConnSkipMark,
		),

		// Skip things from host ip - these are usually kubectl probes
		// skip anything with skip mark. This can be used to add features like port exclusions
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-m", "mark",
			"--mark", constants.SkipMark,
			"-j", "RETURN",
		),

		// Mark outbound connections to route them to the proxy using ip rules/route tables
		// Per Yuval, interface_prefix can be left off this rule... but we should check this (hard to automate
		// detection).
		newIptableRule(
			constants.TableMangle,
			constants.ChainZTunnelPrerouting,
			"-p", "tcp",
			"-m", "set",
			"--match-set", Ipset.Name, "src",
			"-j", "MARK",
			"--set-mark", constants.OutboundMark,
		),
	}

	err = s.iptablesAppend(appendRules)
	if err != nil {
		log.Errorf("failed to append iptables rule: %v", err)
	}

	err = s.iptablesAppend(appendRules2)
	if err != nil {
		log.Errorf("failed to append iptables rule: %v", err)
	}

	// Need to do some work in procfs
	// @TODO: This likely needs to be cleaned up, there are a lot of martians in AWS
	// that seem to necessitate this work.
	procs := map[string]int{
		"/proc/sys/net/ipv4/conf/default/rp_filter":                0,
		"/proc/sys/net/ipv4/conf/all/rp_filter":                    0,
		"/proc/sys/net/ipv4/conf/" + ztunnelVeth + "/rp_filter":    0,
		"/proc/sys/net/ipv4/conf/" + ztunnelVeth + "/accept_local": 1,
	}
	for proc, val := range procs {
		err = SetProc(proc, fmt.Sprint(val))
		if err != nil {
			log.Errorf("failed to write to proc file %s: %v", proc, err)
		}
	}

	// Create tunnels
	inbnd := &netlink.Geneve{
		LinkAttrs: netlink.LinkAttrs{
			Name: constants.InboundTun,
		},
		ID:     1000,
		Remote: net.ParseIP(ztunnelIP),
	}
	log.Debugf("Building inbound tunnel: %+v", inbnd)
	err = netlink.LinkAdd(inbnd)
	if err != nil {
		log.Errorf("failed to add inbound tunnel: %v", err)
	}
	err = netlink.AddrAdd(inbnd, &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   net.ParseIP(constants.InboundTunIP),
			Mask: net.CIDRMask(constants.TunPrefix, 32),
		},
	})
	if err != nil {
		log.Errorf("failed to add inbound tunnel address: %v", err)
	}

	outbnd := &netlink.Geneve{
		LinkAttrs: netlink.LinkAttrs{
			Name: constants.OutboundTun,
		},
		ID:     1001,
		Remote: net.ParseIP(ztunnelIP),
	}
	log.Debugf("Building outbound tunnel: %+v", outbnd)
	err = netlink.LinkAdd(outbnd)
	if err != nil {
		log.Errorf("failed to add outbound tunnel: %v", err)
	}
	err = netlink.AddrAdd(outbnd, &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   net.ParseIP(constants.OutboundTunIP),
			Mask: net.CIDRMask(constants.TunPrefix, 32),
		},
	})
	if err != nil {
		log.Errorf("failed to add outbound tunnel address: %v", err)
	}

	err = netlink.LinkSetUp(inbnd)
	if err != nil {
		log.Errorf("failed to set inbound tunnel up: %v", err)
	}
	err = netlink.LinkSetUp(outbnd)
	if err != nil {
		log.Errorf("failed to set outbound tunnel up: %v", err)
	}

	procs = map[string]int{
		"/proc/sys/net/ipv4/conf/" + constants.InboundTun + "/rp_filter":     0,
		"/proc/sys/net/ipv4/conf/" + constants.InboundTun + "/accept_local":  1,
		"/proc/sys/net/ipv4/conf/" + constants.OutboundTun + "/rp_filter":    0,
		"/proc/sys/net/ipv4/conf/" + constants.OutboundTun + "/accept_local": 1,
	}
	for proc, val := range procs {
		err = SetProc(proc, fmt.Sprint(val))
		if err != nil {
			log.Errorf("failed to write to proc file %s: %v", proc, err)
		}
	}

	dirEntries, err := os.ReadDir("/proc/sys/net/ipv4/conf")
	if err != nil {
		log.Errorf("failed to read /proc/sys/net/ipv4/conf: %v", err)
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			if _, err := os.Stat("/proc/sys/net/ipv4/conf/" + dirEntry.Name() + "/rp_filter"); err != nil {
				err := SetProc("/proc/sys/net/ipv4/conf/"+dirEntry.Name()+"/rp_filter", "0")
				if err != nil {
					log.Errorf("failed to set /proc/sys/net/ipv4/conf/%s/rp_filter: %v", dirEntry.Name(), err)
				}
			}
		}
	}

	routes := []*ExecList{
		newExec("ip",
			[]string{
				"route", "add", "table", fmt.Sprint(constants.RouteTableOutbound), ztunnelIP,
				"dev", ztunnelVeth, "scope", "link",
			},
		),
		newExec("ip",
			[]string{
				"route", "add", "table", fmt.Sprint(constants.RouteTableOutbound), "0.0.0.0/0",
				"via", constants.ZTunnelOutboundTunIP, "dev", constants.OutboundTun,
			},
		),
		newExec("ip",
			[]string{
				"route", "add", "table", fmt.Sprint(constants.RouteTableProxy), ztunnelIP,
				"dev", ztunnelVeth, "scope", "link",
			},
		),
		newExec("ip",
			[]string{
				"route", "add", "table", fmt.Sprint(constants.RouteTableProxy), "0.0.0.0/0",
				"via", ztunnelIP, "dev", ztunnelVeth, "onlink",
			},
		),
		newExec("ip",
			[]string{
				"route", "add", "table", fmt.Sprint(constants.RouteTableInbound), ztunnelIP,
				"dev", ztunnelVeth, "scope", "link",
			},
		),
		// Everything with the skip mark goes directly to the main table
		newExec("ip",
			[]string{
				"rule", "add", "priority", "100",
				"fwmark", fmt.Sprint(constants.SkipMark),
				"goto", "32766",
			},
		),
		// Everything with the outbound mark goes to the tunnel out device
		// using the outbound route table
		newExec("ip",
			[]string{
				"rule", "add", "priority", "101",
				"fwmark", fmt.Sprint(constants.OutboundMark),
				"lookup", fmt.Sprint(constants.RouteTableOutbound),
			},
		),
		// Things with the proxy return mark go directly to the proxy veth using the proxy
		// route table (useful for original src)
		newExec("ip",
			[]string{
				"rule", "add", "priority", "102",
				"fwmark", fmt.Sprint(constants.ProxyRetMark),
				"lookup", fmt.Sprint(constants.RouteTableProxy),
			},
		),
		// Send all traffic to the inbound table. This table has routes only to pods in the mesh.
		// It does not have a catch-all route, so if a route is missing, the search will continue
		// allowing us to override routing just for member pods.
		newExec("ip",
			[]string{
				"rule", "add", "priority", "103",
				"table", fmt.Sprint(constants.RouteTableInbound),
			},
		),
	}

	for _, route := range routes {
		err = execute(route.Cmd, route.Args...)
		if err != nil {
			log.Errorf(fmt.Errorf("failed to add route (%+v): %v", route, err))
		}
	}

	return nil
}

func (s *Server) updateZtunnelEbpfOnNode(pod *corev1.Pod, captureDNS bool) error {
	if s.ebpfServer == nil {
		return fmt.Errorf("uninitialized ebpf server")
	}

	ip := pod.Status.PodIP
	args, err := buildEbpfArgsByIP(ip, true, false)
	if err != nil {
		return err
	}
	args.CaptureDNS = captureDNS

	log.Debugf("update ztunnel ebpf args: %+v", args)
	s.ebpfServer.AcceptRequest(args)
	return nil
}

func (s *Server) delZtunnelEbpfOnNode() error {
	if s.ebpfServer == nil {
		return fmt.Errorf("uninitialized ebpf server")
	}

	args := &ebpf.RedirectArgs{
		Ifindex:   0,
		IsZtunnel: true,
		Remove:    true,
	}
	log.Debugf("del ztunnel ebpf args: %+v", args)
	s.ebpfServer.AcceptRequest(args)
	return nil
}

func (s *Server) updatePodEbpfOnNode(pod *corev1.Pod) error {
	if s.ebpfServer == nil {
		return fmt.Errorf("uninitialized ebpf server")
	}

	ip := pod.Status.PodIP
	if ip == "" {
		log.Debugf("skip adding pod %s/%s, IP not yet allocated", pod.Name, pod.Namespace)
		return nil
	}

	args, err := buildEbpfArgsByIP(ip, false, false)
	if err != nil {
		return err
	}

	log.Debugf("update POD ebpf args: %+v", args)
	s.ebpfServer.AcceptRequest(args)
	return nil
}

func (s *Server) AddPodToMesh(pod *corev1.Pod) {
	switch s.redirectMode {
	case IptablesMode:
		addPodToMeshWithIptables(s.kubeClient.Kube(), pod, "")
	case EbpfMode:
		if err := s.updatePodEbpfOnNode(pod); err != nil {
			log.Errorf("failed to update POD ebpf: %v", err)
		}
		if err := AnnotateEnrolledPod(s.kubeClient.Kube(), pod); err != nil {
			log.Errorf("failed to annotate pod enrollment: %v", err)
		}
	}
}

func (s *Server) delPodEbpfOnNode(ip string) error {
	if s.ebpfServer == nil {
		return fmt.Errorf("uninitialized ebpf server")
	}

	if ip == "" {
		log.Debugf("nothing could be performed to delete ebpf for empty ip")
		return nil
	}
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return fmt.Errorf("failed to parse ip(%s): %v", ip, err)
	}

	ifIndex := 0

	if veth, err := getVethWithDestinationOf(ip); err != nil {
		log.Infof("failed to get device: %v", err) // Change to debug, it's expected intf has been destroyed
	} else {
		ifIndex = veth.Attrs().Index
	}

	args := &ebpf.RedirectArgs{
		IPAddrs:   []netip.Addr{ipAddr},
		Ifindex:   ifIndex,
		IsZtunnel: false,
		Remove:    true,
	}
	log.Debugf("del POD ebpf args: %+v", args)
	s.ebpfServer.AcceptRequest(args)
	return nil
}

func (s *Server) DelPodFromMesh(pod *corev1.Pod) {
	switch s.redirectMode {
	case IptablesMode:
		DelPodFromMesh(s.kubeClient.Kube(), pod)
	case EbpfMode:
		if pod.Spec.HostNetwork {
			log.Debugf("pod(%s/%s) is using host network, skip it", pod.Namespace, pod.Name)
			return
		}
		if err := s.delPodEbpfOnNode(pod.Status.PodIP); err != nil {
			log.Errorf("failed to del POD ebpf: %v", err)
		}
		if err := AnnotateUnenrollPod(s.kubeClient.Kube(), pod); err != nil {
			log.Errorf("failed to annotate pod unenrollment: %v", err)
		}
	}
}

func (s *Server) cleanup() {
	log.Infof("server terminated, cleaning up")
	if s.redirectMode == EbpfMode {
		if err := s.delZtunnelEbpfOnNode(); err != nil {
			log.Error(err)
		}
		return
	}
	s.cleanRules()

	// Clean up ip route tables
	_ = routeFlushTable(constants.RouteTableInbound)
	_ = routeFlushTable(constants.RouteTableOutbound)
	_ = routeFlushTable(constants.RouteTableProxy)

	exec := []*ExecList{
		newExec("ip", []string{"rule", "del", "priority", "100"}),
		newExec("ip", []string{"rule", "del", "priority", "101"}),
		newExec("ip", []string{"rule", "del", "priority", "102"}),
		newExec("ip", []string{"rule", "del", "priority", "103"}),
	}
	for _, e := range exec {
		err := execute(e.Cmd, e.Args...)
		if err != nil {
			log.Warnf("Error running command %v %v: %v", e.Cmd, strings.Join(e.Args, " "), err)
		}
	}

	// Delete tunnel links
	err := netlink.LinkDel(&netlink.Geneve{
		LinkAttrs: netlink.LinkAttrs{
			Name: constants.InboundTun,
		},
	})
	if err != nil {
		log.Warnf("error deleting inbound tunnel: %v", err)
	}
	err = netlink.LinkDel(&netlink.Geneve{
		LinkAttrs: netlink.LinkAttrs{
			Name: constants.OutboundTun,
		},
	})
	if err != nil {
		log.Warnf("error deleting outbound tunnel: %v", err)
	}

	_ = Ipset.DestroySet()
}

func routeFlushTable(table int) error {
	routes, err := netlink.RouteListFiltered(netlink.FAMILY_V4, &netlink.Route{Table: table}, netlink.RT_FILTER_TABLE)
	if err != nil {
		return err
	}
	err = routesDelete(routes)
	if err != nil {
		return err
	}
	return nil
}

func routesDelete(routes []netlink.Route) error {
	for _, r := range routes {
		err := netlink.RouteDel(&r)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetProc(path string, value string) error {
	return os.WriteFile(path, []byte(value), 0o644)
}
