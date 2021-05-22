package src

import (
	"sync"
)

type NodeSet struct {
	onlineNodes  map[NodeAddress]bool
	onlineNodeMX sync.RWMutex
}

func NewNodeSet() NodeSet {
	return NodeSet{
		onlineNodes:  map[NodeAddress]bool{},
		onlineNodeMX: sync.RWMutex{},
	}
}

func (set *NodeSet) AddOnlineNode(node NodeAddress) {
	set.AddOnlineNodes([]NodeAddress{node})
}
func (set *NodeSet) AddOnlineNodes(nodes []NodeAddress) {
	set.onlineNodeMX.Lock()
	defer set.onlineNodeMX.Unlock()

	for _, v := range nodes {
		set.onlineNodes[v] = true
	}

}

func (set *NodeSet) GetOnlineNodes() []NodeAddress {
	set.onlineNodeMX.RLock()
	defer set.onlineNodeMX.RUnlock()

	// convert set set to list
	onlineNodeList := []NodeAddress{}
	for k, _ := range set.onlineNodes {
		onlineNodeList = append(onlineNodeList, k)
	}

	return onlineNodeList
}

func (set *NodeSet) RemoveOnlineNode(nodes []NodeAddress) {
	set.onlineNodeMX.Lock()
	defer set.onlineNodeMX.Unlock()

	for _, node := range nodes {
		delete(set.onlineNodes, node)
	}
}

// return a list of nodes in diffNodes that are not present in set
// this.set - setB
// eg. diffNodes = [node1, node 2, node3], set=[node 1, node 2]
// GetSetDiff(diffNodes) => [node3]
func (set NodeSet) GetDifference(setB []NodeAddress) []NodeAddress {
	set.onlineNodeMX.RLock()
	defer set.onlineNodeMX.RUnlock()

	diff := []NodeAddress{}

	// go through all provided nodes
	for _, n := range setB {
		if _, ok := set.onlineNodes[n]; !ok {
			// if not found in set add to list
			diff = append(diff, n)
		}
	}

	return diff
}
