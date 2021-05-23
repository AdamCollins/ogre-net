package node

import (
	"github.com/AdamCollins/ogre-net/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStart(t *testing.T) {

	// Node starts up knowing 1 existing node
	node := Node{}
	conHandler := MockConnectionHandler{}
	config := Config{
		NodeId:     "TestNode",
		KnownNodes: []types.NodeAddress{":3000"},
		ListenAddr: ":3001",
	}

	StartWithHandler(config, &conHandler, &node)

	// test node was set up
	assert.Equal(t, &conHandler, node.conHandler, "node's connection handler should be set")
	assert.Equal(t, config.NodeId, node.NodeId, "node's NodeId should be set")
	assert.Equal(t, node.ListenAddr, config.ListenAddr, node.ListenAddr, "node's ListenAddr should be set")

	// test made correct calls to connectionHandler
	assert.Equal(t, []types.NodeAddress{":3001"}, conHandler.ListenRequestLog)
	assert.Equal(t, []types.NodeAddress{":3000", ":3001"}, conHandler.PingRequestLog)
	assert.Equal(t, []AdRequest{{
		Target:      ":3000",
		Advertising: []types.NodeAddress{":3000", ":3001"},
	}}, conHandler.AdvertiseRequestLog)

}
