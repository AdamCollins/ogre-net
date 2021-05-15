package src

import (
	"log"
)

type NodeConfig struct {
	NodeId     string
	KnownNodes []NodeAddress
	ListenAddr NodeAddress
}

type Node struct {
	NodeId      string
	onlineNodes NodeSet
	handler     INodeRPCHandler
	ListenAddr  NodeAddress
}

func StartListening(config NodeConfig) {
	rpcHandler := NodeRPCHandler{}
	StartListeningWithHandler(config, rpcHandler)
}

func StartListeningWithHandler(config NodeConfig, handler INodeRPCHandler) {
	// set handler to injected handler

	// start listening for RPC requests on config.ListenAddr in new go routine
	// go handler.Listen()

	// Add knownNodes to onlineNodes

	// Add ListenAddr to OnlineNodes

	// Contact nodes in config.KnownNodes
	//PingNodes()

}

func (node Node) PingOnlineNodes() {

	newNeighbourSet := map[NodeAddress]bool{}

	// get list of all known online nodes
	currentOnlineNodes := node.onlineNodes.GetOnlineNodes()

	// ping all known online nodes and ask for neighbours
	for _, nodeAddr := range currentOnlineNodes {
		jlkjlkjlkjklkkkkk
		// get a diff list of neighbours

		// no need to call self. skip
		if nodeAddr == node.ListenAddr {
			continue
		}

		// maybe put in go routine?
		neighbours, err := node.handler.GetNodesNeighbours(nodeAddr, currentOnlineNodes)
		if err != nil {
			// if neighbour does not respond remove it from list of nodes.
			log.Printf("[Node]: %v could not find node at %v. Removing from list\n", node.NodeId, nodeAddr)
			node.onlineNodes.RemoveOnlineNode(nodeAddr)
		}

		for _, neighbour := range neighbours {
			if _, ok := newNeighbourSet[neighbour]; ok {
				newNeighbourSet[neighbour] = true
			}
		}
	}

	newNeighbourList := make([]NodeAddress, len(newNeighbourSet))
	for newNeighbour, _ := range newNeighbourSet {

		// verify that the received node is actually online.
		err := node.handler.PingNode(newNeighbour)
		if err == nil {
			newNeighbourList = append(newNeighbourList, newNeighbour)
			log.Printf("[Node]: %v is adding node at %v\n", node, newNeighbour)
		}
	}

	// Add new nodes to online Nodes
	node.onlineNodes.AddOnlineNodes(newNeighbourList)

	log.Printf("[Node]: %v just finished adding %d neightbours\n", node, len(newNeighbourList))
}

func (node Node) GetOnlineNodesDiff(callerNodes []NodeAddress) []NodeAddress {
	return node.onlineNodes.GetSetDiff(callerNodes)
}
