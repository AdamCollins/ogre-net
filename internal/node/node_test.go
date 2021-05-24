package node

import (
	"github.com/AdamCollins/ogre-net/internal/types"
	"github.com/AdamCollins/ogre-net/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tests = []Config{
	{
		NodeId:     "SingleOnlineNode",
		KnownNodes: []types.NodeAddress{":3000"},
		ListenAddr: ":3001",
	},
	{
		NodeId:     "NoOnlineNodes",
		KnownNodes: []types.NodeAddress{},
		ListenAddr: ":3001",
	},
	{
		NodeId:     "ThreeOnlineNodes",
		KnownNodes: []types.NodeAddress{":3002", ":3003", ":3004"},
		ListenAddr: ":3001",
	},
}

func TestStart_NodesOnline(t *testing.T) {

	for _, test := range tests {
		t.Run(test.NodeId, func(t *testing.T) {
			node := Node{}
			// Creates a mock connection handler in which all requested contacted nodes respond to RPC calls
			conHandler := mocks.MockConnectionHandlerAllOnline{}
			StartWithHandler(test, &conHandler, &node)

			// test node was set up
			assert.Equal(t, &conHandler, node.conHandler, "node's connection handler should be set")
			assert.Equal(t, test.NodeId, node.NodeId, "node's NodeId should be set")
			assert.Equal(t, node.ListenAddr, test.ListenAddr, node.ListenAddr, "node's ListenAddr should be set")

			// test made correct calls to connectionHandler
			assert.Equal(t, []types.NodeAddress{test.ListenAddr}, conHandler.ListenRequestLog)
			assert.Equal(t, append(test.KnownNodes, test.ListenAddr), conHandler.PingRequestLog)

			for i, onlineNode := range test.KnownNodes {
				assert.Equal(t, mocks.AdRequest{
					Target:      onlineNode,
					Advertising: append(test.KnownNodes, test.ListenAddr),
				}, conHandler.AdvertiseRequestLog[i])
			}
		})
	}
}

func TestStart_NoOnlineNodes(t *testing.T) {

	for _, test := range tests {
		t.Run(test.NodeId, func(t *testing.T) {
			node := Node{}
			// Creates a mock connection handler in which all requested contacted nodes fail to respond
			conHandler := mocks.NewMockOffline(test.ListenAddr)
			StartWithHandler(test, &conHandler, &node)

			// test node was set up
			assert.Equal(t, &conHandler, node.conHandler, "node's connection handler should be set")
			assert.Equal(t, test.NodeId, node.NodeId, "node's NodeId should be set")
			assert.Equal(t, node.ListenAddr, test.ListenAddr, node.ListenAddr, "node's ListenAddr should be set")

			// test made correct calls to connectionHandler
			assert.Equal(t, []types.NodeAddress{test.ListenAddr}, conHandler.ListenRequestLog)
			// should only send ping to self
			assert.Equal(t, []types.NodeAddress{test.ListenAddr}, conHandler.PingRequestLog)
			// no advertisements should be sent out
			assert.Equal(t, []mocks.AdRequest{}, conHandler.AdvertiseRequestLog)
		})
	}
}
