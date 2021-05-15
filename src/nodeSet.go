package src

import "sync"

type NodeSet struct {
	onlineNodes  map[NodeAddress]bool
	onlineNodeMX sync.RWMutex
}

func (set NodeSet) AddOnlineNodes(nodes []NodeAddress) {
	set.onlineNodeMX.Lock()
	defer set.onlineNodeMX.Unlock()

	for _, v := range nodes {
		set.onlineNodes[v] = true
	}

}

func (set NodeSet) GetOnlineNodes() []NodeAddress {
	set.onlineNodeMX.RLock()
	defer set.onlineNodeMX.RUnlock()

	// convert set set to list
	onlineNodeList := make([]NodeAddress, len(set.onlineNodes))
	for k, _ := range set.onlineNodes {
		onlineNodeList = append(onlineNodeList, k)
	}

	return onlineNodeList
}

func (set NodeSet) RemoveOnlineNode(address NodeAddress) {
	set.onlineNodeMX.Lock()
	defer set.onlineNodeMX.Unlock()

	delete(set.onlineNodes, address)

}

// return a list of nodes in diffNodes that are not present in set
func (set NodeSet) GetSetDiff(diffNodes []NodeAddress) []NodeAddress {
	set.onlineNodeMX.RLock()
	defer set.onlineNodeMX.RUnlock()

	var diff []NodeAddress

	// go through all provided nodes
	for _, n := range diffNodes {
		if _, ok := set.onlineNodes[n]; !ok {
			// if not found in set add to list
			diff = append(diff, n)
		}
	}

	return diff
}
