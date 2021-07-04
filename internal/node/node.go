package node

import (
	"github.com/AdamCollins/ogre-net/internal/node/node_set"
	"github.com/AdamCollins/ogre-net/internal/onion"
	"github.com/AdamCollins/ogre-net/internal/types"
	"github.com/AdamCollins/ogre-net/internal/utils"
	"log"
	"sync"
	"time"
)

var Constants = types.Constant

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

	// start Listening for requests
	StartWithHandler(config, &rpcHandler, &node)

	ticker := time.NewTicker(time.Duration(Constants.PingInterval) * time.Millisecond)
	// continuously check that nodes are still online
	for {
		select {
		case <-ticker.C:
			//node.PingAllOnlineNodes()
		}
	}
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

// receives a list of potential new nodes
// ensure that they are all online and then sends out advertisement if any new nodes are found.
func (node *Node) AddOnlineNodes(newNodes []types.NodeAddress) {

	// first make sure node is really online
	newOnlineNodes, _ := node.pingNodes(newNodes)

	// next add to onlineNode set
	node.onlineNodes.AddOnlineNodes(newOnlineNodes)
	currentOnlineNodes := node.onlineNodes.GetNodes()
	log.Printf("[Node: %v] onlineNodes now include %v\n", node.NodeId, currentOnlineNodes)

	// now advertise to neighbours
	node.Advertise()
}

// sends advertisement of onlineNodes to all onlineNodes(excluding self)
func (node *Node) Advertise() {

	// store any nodes that have died
	var deadNodes []types.NodeAddress

	// get list of all known online nodes
	currentOnlineNodes := node.onlineNodes.GetNodes()
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
		node.onlineNodes.RemoveOnlineNodes(deadNodes)
	}

}

// checks to see if there are any callerNodes that this node does not already know about.
// if so, add these nodes
func (node *Node) ReceiveAdvertisement(callerNodes []types.NodeAddress) {
	log.Printf("[Node: %v] received advertisment with nodes %v\n", node.NodeId, callerNodes)
	newNodes := node.onlineNodes.GetDifference(callerNodes)
	if len(newNodes) > 0 {
		log.Printf("[Node: %v] Adding nodes %v\n", node.NodeId, newNodes)
		node.AddOnlineNodes(newNodes)
	}
}

func (node *Node) PingAllOnlineNodes() {
	// ping all online nodes
	_, died := node.pingNodes(node.onlineNodes.GetNodes())
	// remove any that have died
	if len(died) > 0 {
		node.onlineNodes.RemoveOnlineNodes(died)
	}
}

// returns a random slice of online nodes with a max length of n
func (node *Node) GetRandomNodesSubset(n uint16) []types.NodeAddress {
	nodes := node.onlineNodes.GetNodes()

	nodes = utils.ShuffleNodes(nodes)

	// return no more than n nodes
	maxLen := utils.MinUInt16(uint16(len(nodes)), n)
	return nodes[:maxLen]
}

// takes an OnionMessage, msg, and peels it by one layer.
// If no more hops need to be made(NextLayer==nil) make request in payload.
// 		- take the response from the payload and wrap it in an OnionMessage
// Otherwise, forward to NextHop
// 		- take response from the forward and add a layer to the message before responding with it.
func (node *Node) ReceiveMessage(msg types.OnionMessage) types.OnionMessage {
	log.Printf("[Node %v] received message: %v \n", node.NodeId, msg)
	peeledMsg := onion.Peel(msg)
	log.Printf("[Node %v] peeled out message: %v\n", node.NodeId, peeledMsg)

	var response types.OnionMessage
	if peeledMsg.NextLayer == nil {
		// done. this is the exit node
		log.Printf("[Node %v] received payload '%v'\n", node.NodeId, peeledMsg.Payload)
		// handle request in payload
		responseText := handleExitRequest(peeledMsg.Payload)
		// create new onionMessage for response
		response = onion.NewOnionMessage(responseText, []types.NodeAddress{node.ListenAddr})
	} else {
		// forward
		log.Printf("[Node %v] forwarding message to %v'\n", node.NodeId, peeledMsg.NextHop)
		response, _ = node.conHandler.ForwardMessage(peeledMsg, peeledMsg.NextHop)
		// get response and self as next layer
		response = onion.Layer(response, node.ListenAddr)
	}

	return response
}

func handleExitRequest(payload string) string {
	//TODO
	return "Made request " + payload
}

// pings provided nodes and returns a list of alive nodes and a list of dead nodes
func (node *Node) pingNodes(nodes []types.NodeAddress) (alive []types.NodeAddress, dead []types.NodeAddress) {
	// channel to store any dead channels
	// (performance? we can assume that most pinged nodes will be alive so might be worth reducing chan size
	// to reduce malloc time? blocking on full channel in this case would be alright)
	var deadNodes = make(chan types.NodeAddress, len(nodes))
	var aliveNodes = make(chan types.NodeAddress, len(nodes))
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for _, nodeAddr := range nodes {
		go func(na types.NodeAddress) {
			defer wg.Done()
			err := node.conHandler.PingNode(na)
			if err != nil {
				// if no response add to deadNode
				deadNodes <- na
				log.Printf("[Node: %v]:No response from node at %v. Removing from list\n", node.NodeId, na)
			} else {
				// if response is received add to aliveNodes
				aliveNodes <- na
			}
		}(nodeAddr)
	}
	// wait for all requests to be complete
	wg.Wait()
	close(aliveNodes)
	close(deadNodes)
	return utils.ChanToSlice(aliveNodes), utils.ChanToSlice(deadNodes)
}
