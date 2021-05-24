package node

import (
	"github.com/AdamCollins/ogre-net/internal/types"
	"github.com/AdamCollins/ogre-net/mocks"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type TestResult struct {
	ListenRequests []types.NodeAddress
	PingResponses  []types.NodeAddress
	AdRequests     []mocks.AdRequest
}

type Test struct {
	config   Config
	expected TestResult
}

var testsNodesOnline = []Test{
	{
		Config{
			NodeId:     "SingleOnlineNode",
			KnownNodes: []types.NodeAddress{":3000"},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3000", ":3001"},
			[]mocks.AdRequest{
				{
					types.NodeAddress(":3000"),
					[]types.NodeAddress{":3000", ":3001"},
				},
			},
		},
	},
	{
		Config{
			NodeId:     "NoOnlineNodes",
			KnownNodes: []types.NodeAddress{},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3001"},
			[]mocks.AdRequest{},
		},
	},
	{
		Config{
			NodeId:     "ThreeOnlineNodes",
			KnownNodes: []types.NodeAddress{":3002", ":3003", ":3004"},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3001", ":3002", ":3003", ":3004"},
			[]mocks.AdRequest{
				{
					types.NodeAddress(":3002"),
					[]types.NodeAddress{":3001", ":3002", ":3003", ":3004"},
				},
				{
					types.NodeAddress(":3003"),
					[]types.NodeAddress{":3001", ":3002", ":3003", ":3004"},
				},
				{
					types.NodeAddress(":3004"),
					[]types.NodeAddress{":3001", ":3002", ":3003", ":3004"},
				},
			},
		},
	},
}

func TestStart_NodesOnline(t *testing.T) {

	for _, test := range testsNodesOnline {
		t.Run(test.config.NodeId, func(t *testing.T) {
			node := Node{}
			// Creates a mock connection handler in which all requested contacted nodes respond to RPC calls
			conHandler := mocks.MockConnectionHandlerAllOnline{}
			StartWithHandler(test.config, &conHandler, &node)

			// test node was set up
			assert.Equal(t, &conHandler, node.conHandler, "node's connection handler should be set")
			assert.Equal(t, test.config.NodeId, node.NodeId, "node's NodeId should be set")
			assert.Equal(t, test.config.ListenAddr, node.ListenAddr, "node's ListenAddr should be set")

			// test made correct calls to connectionHandler
			assert.True(t, expectedNodeAddrSlice(test.expected.ListenRequests, conHandler.ListenRequestLog), "Listen request sent")

			// should ping itself and any other online nodes
			assert.True(t, expectedNodeAddrSlice(test.expected.PingResponses, conHandler.PingRequestLog))

			// should advertise to any online nodes but not to itself
			assert.True(t, expectedAdRequests(test.expected.AdRequests, conHandler.AdvertiseRequestLog))
		})
	}
}

var testsNoNodesOnline = []Test{
	{
		Config{
			NodeId:     "SingleOnlineNode",
			KnownNodes: []types.NodeAddress{":3000"},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3001"},
			[]mocks.AdRequest{},
		},
	},
	{
		Config{
			NodeId:     "NoOnlineNodes",
			KnownNodes: []types.NodeAddress{},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3001"},
			[]mocks.AdRequest{},
		},
	},
	{
		Config{
			NodeId:     "ThreeOnlineNodes",
			KnownNodes: []types.NodeAddress{":3002", ":3003", ":3004"},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3001"},
			[]mocks.AdRequest{},
		},
	},
}

func TestStart_NoOnlineNodes(t *testing.T) {

	for _, test := range testsNoNodesOnline {
		t.Run(test.config.NodeId, func(t *testing.T) {
			node := Node{}
			// Creates a mock connection handler in which all requested contacted nodes fail to respond
			conHandler := mocks.NewMockOffline(test.config.ListenAddr)
			StartWithHandler(test.config, &conHandler, &node)

			// test made correct calls to connectionHandler
			assert.True(t, expectedNodeAddrSlice(test.expected.ListenRequests, conHandler.ListenRequestLog), "Listen request sent")

			// should ping itself and any other online nodes
			assert.True(t, expectedNodeAddrSlice(test.expected.PingResponses, conHandler.PingRequestLog))

			// should advertise to any online nodes but not to itself
			assert.True(t, expectedAdRequests(test.expected.AdRequests, conHandler.AdvertiseRequestLog))
		})
	}
}

var testsOneNodeOnline = []Test{
	{
		Config{
			NodeId:     "SingleOnlineNode",
			KnownNodes: []types.NodeAddress{":3000"},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3001", ":3000"},
			[]mocks.AdRequest{
				{
					":3000",
					[]types.NodeAddress{":3000", ":3001"},
				},
			},
		},
	},
	{
		Config{
			NodeId:     "NoOnlineNodes",
			KnownNodes: []types.NodeAddress{},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3001"},
			[]mocks.AdRequest{},
		},
	},
	{
		Config{
			NodeId:     "ThreeOnlineNodes",
			KnownNodes: []types.NodeAddress{":3002", ":3003", ":3004"},
			ListenAddr: ":3001",
		},
		TestResult{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3001"},
			[]mocks.AdRequest{},
		},
	},
}

func TestStart_OneOnlineNode(t *testing.T) {

	for _, test := range testsOneNodeOnline {
		t.Run(test.config.NodeId, func(t *testing.T) {
			node := Node{}
			// Creates a mock connection handler in which all requested contacted nodes fail to respond
			conHandler := mocks.NewMockOneOnline(test.config.ListenAddr)
			StartWithHandler(test.config, &conHandler, &node)

			// test made correct calls to connectionHandler
			assert.True(t, expectedNodeAddrSlice(test.expected.ListenRequests, conHandler.ListenRequestLog), "Listen request sent")

			// should ping itself and any other online nodes
			assert.True(t, expectedNodeAddrSlice(test.expected.PingResponses, conHandler.PingResponseLog))

			// should advertise to any online nodes but not to itself
			assert.True(t, expectedAdRequests(test.expected.AdRequests, conHandler.AdvertiseResponseLog))
		})
	}
}

//func TestNode_ReceiveAdvertisementNewNode(t *testing.T) {
//
//	conHandler := mocks.MockConnectionHandlerAllOnline{}
//	onlineNodes := node_set.NewNodeSet()
//	onlineNodes.AddOnlineNode(":3001")
//	// node with no onlineNodes
//	node := Node{
//		onlineNodes: onlineNodes,
//		conHandler:  &conHandler,
//		ListenAddr:  ":3001",
//	}
//
//	ad := []types.NodeAddress{":3002"}
//	// send node an advertisement with a new node
//	node.ReceiveAdvertisement(ad)
//
//	// node should ping the new node
//	assert.Equal(t, []types.NodeAddress{":3002"}, conHandler.PingRequestLog)
//	// node should add the node to its onlineNodes set
//	assert.Equal(t, 2, len(node.onlineNodes.GetNodes()))
//	// node should send advertisement to node
//	assert.Equal(t, 1, len(conHandler.AdvertiseRequestLog))
//}func TestNode_ReceiveAdvertisementNewNode(t *testing.T) {
//
//	conHandler := mocks.MockConnectionHandlerAllOnline{}
//	onlineNodes := node_set.NewNodeSet()
//	onlineNodes.AddOnlineNode(":3001")
//	// node with no onlineNodes
//	node := Node{
//		onlineNodes: onlineNodes,
//		conHandler:  &conHandler,
//		ListenAddr:  ":3001",
//	}
//
//	ad := []types.NodeAddress{":3002"}
//	// send node an advertisement with a new node
//	node.ReceiveAdvertisement(ad)
//
//	// node should ping the new node
//	assert.Equal(t, []types.NodeAddress{":3002"}, conHandler.PingRequestLog)
//	// node should add the node to its onlineNodes set
//	assert.Equal(t, 2, len(node.onlineNodes.GetNodes()))
//	// node should send advertisement to node
//	assert.Equal(t, 1, len(conHandler.AdvertiseRequestLog))
//}

// test utils

func expectedAdRequests(expected, actual []mocks.AdRequest) bool {
	if len(expected) != len(actual) {
		log.Printf("was expecting %v but got %v\n", expected, actual)
		return false
	}
	dict := map[types.NodeAddress][]types.NodeAddress{}

	for _, adA := range actual {
		dict[adA.Target] = adA.Advertisement
	}

	for _, adB := range expected {
		doesAdExists := dict[adB.Target] != nil

		if !doesAdExists {
			log.Printf("was expecting %v but it was not found\n", adB.Target)
			log.Printf("was expecting %v but got %v\n", expected, actual)
			return false
		}

		isAdIdentical := expectedNodeAddrSlice(dict[adB.Target], adB.Advertisement)
		if !isAdIdentical {
			log.Printf("was expecting %v but got %v\n", expected, actual)
			return false
		}
	}

	return true

}

// returns true  if sliceA and sliceB are equal ignoring order
func expectedNodeAddrSlice(expected, actual []types.NodeAddress) bool {
	if len(expected) != len(actual) {
		log.Printf("was expecting %v but got %v\n", expected, actual)
		return false
	}

	dict := map[types.NodeAddress]bool{}
	for _, itemA := range expected {
		dict[itemA] = true
	}

	for _, itemB := range actual {
		if !dict[itemB] {
			log.Printf("was expecting %v but got %v\n", expected, actual)
			return false
		}
	}

	return true
}
