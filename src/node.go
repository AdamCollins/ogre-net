package src

import (
	"log"
	"sync"
)

type NodeConfig struct {
	NodeId     string
	KnownNodes []NodeAddress
	ListenAddr NodeAddress
}

type Node struct {
	NodeId      string
	onlineNodes NodeSet
	connection  INodeConnectionHandler
	ListenAddr  NodeAddress
}

func StartListening(config NodeConfig) {
	rpcHandler := NodeConnectionHandler{}
	StartListeningWithHandler(config, &rpcHandler)
}

func StartListeningWithHandler(config NodeConfig, handler INodeConnectionHandler) {
	// set connection to injected connection
	node := Node{
		NodeId:      config.NodeId,
		onlineNodes: NewNodeSet(),
		connection:  handler,
		ListenAddr:  config.ListenAddr,
	}

	// start listening for RPC requests on config.ListenAddr in new go routine
	err := handler.Listen(config.ListenAddr, &node)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Node: %v] now listening for connections at %v\n", node.NodeId, node.ListenAddr)

	// start advertising for known nodes and then wait for advertisements back
	initNodes := append(config.KnownNodes, config.ListenAddr)
	node.AddOnlineNodes(initNodes)

	// forever yield while waiting for requests\
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func (node *Node) AddOnlineNodes(newNodes []NodeAddress) {

	var newOnlineNodes []NodeAddress
	// first make sure node is really online
	for _, newNode := range newNodes {
		err := node.connection.PingNode(newNode)
		if err != nil {
			// node isn't online
			log.Printf("[Node: %v] tried to add node %v but was found to be offline or non-existent.\n", node.NodeId, newNode)
			continue
		}
		newOnlineNodes = append(newOnlineNodes, newNode)
	}
	// next add to onlineNode set
	node.onlineNodes.AddOnlineNodes(newOnlineNodes)
	currentOnlineNodes := node.onlineNodes.GetOnlineNodes()
	log.Printf("[Node: %v] onlineNodes now include %v\n", node.NodeId, currentOnlineNodes)

	// now advertise to neighbours
	node.Advertise()
}

func (node *Node) Advertise() {

	// store any nodes that have died
	var deadNodes []NodeAddress

	// get list of all known online nodes
	currentOnlineNodes := node.onlineNodes.GetOnlineNodes()
	// advertise known nodes to all known online nodes
	for _, nodeAddr := range currentOnlineNodes {

		// no need to call self. skip
		if nodeAddr == node.ListenAddr {
			continue
		}

		// get a diff list of neighbours
		// maybe put in go routine?
		err := node.connection.Advertise(nodeAddr, currentOnlineNodes)
		if err != nil {
			// if neighbour does not respond remove it from list of nodes.
			log.Printf("[Node: %v]:could not find node at %v. Removing from list\n", node.NodeId, nodeAddr)
			deadNodes = append(deadNodes, nodeAddr)
		}
	}

	// if any nodes have died remove them
	if deadNodes != nil {
		node.onlineNodes.RemoveOnlineNode(deadNodes)
	}

}

// checks to see if there are any callerNodes that this node does not already know about.
// if so, add these nodes
func (node *Node) CheckForNewNodes(callerNodes []NodeAddress) {
	log.Printf("[Node: %v] received advertisment with nodes %v\n", node.NodeId, callerNodes)
	newNodes := node.onlineNodes.GetDifference(callerNodes)
	if len(newNodes) > 0 {
		log.Printf("[Node: %v] Adding nodes %v\n", node.NodeId, newNodes)

		node.AddOnlineNodes(newNodes)
	}
}
