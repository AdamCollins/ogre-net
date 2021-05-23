package node

import (
	"github.com/AdamCollins/ogre-net/internal/node/node_set"
	"github.com/AdamCollins/ogre-net/internal/types"
	"log"
	"sync"
)

type Config struct {
	NodeId     string
	KnownNodes []types.NodeAddress
	ListenAddr types.NodeAddress
}

type Node struct {
	NodeId      string
	onlineNodes node_set.NodeSet
	conHandler  types.INodeConnectionHandler
	ListenAddr  types.NodeAddress
}

func Start(config Config) {
	node := Node{}
	rpcHandler := ConnectionHandler{}
	rpcHandler.SetCallbackNode(&node)
	StartWithHandler(config, &rpcHandler, &node)
	// forever yield while waiting for requests
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func StartWithHandler(config Config, handler types.INodeConnectionHandler, node *Node) {
	// set conHandler to injected conHandler
	*node = Node{
		NodeId:      config.NodeId,
		onlineNodes: node_set.NewNodeSet(),
		conHandler:  handler,
		ListenAddr:  config.ListenAddr,
	}

	// start listening for RPC requests on config.ListenAddr in new go routine
	err := node.conHandler.Listen(config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Node: %v] now listening for connections at %v\n", node.NodeId, node.ListenAddr)

	// start advertising for known nodes and then wait for advertisements back
	initNodes := append(config.KnownNodes, config.ListenAddr)
	node.AddOnlineNodes(initNodes)
}

func (node *Node) AddOnlineNodes(newNodes []types.NodeAddress) {

	var newOnlineNodes []types.NodeAddress
	// first make sure node is really online
	for _, newNode := range newNodes {
		err := node.conHandler.PingNode(newNode)
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
	var deadNodes []types.NodeAddress

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
		err := node.conHandler.Advertise(nodeAddr, currentOnlineNodes)
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
func (node *Node) CheckForNewNodes(callerNodes []types.NodeAddress) {
	log.Printf("[Node: %v] received advertisment with nodes %v\n", node.NodeId, callerNodes)
	newNodes := node.onlineNodes.GetDifference(callerNodes)
	if len(newNodes) > 0 {
		log.Printf("[Node: %v] Adding nodes %v\n", node.NodeId, newNodes)

		node.AddOnlineNodes(newNodes)
	}
}
