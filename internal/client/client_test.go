package client

import (
	"github.com/AdamCollins/ogre-net/internal/types"
	"github.com/AdamCollins/ogre-net/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartHelper_SuperNodeDown(t *testing.T) {
	conHandler := mocks.ClientMockConnectionHandler{}
	t.Run("attempt start with super node offline", func(t *testing.T) {
		config := Config{
			IdealNumberOfHops: 4,
			MinNumberOfHops:   2,
			KnownSuperNode:    ":3001",
		}

		_, err := StartHelper(config, &conHandler)
		assert.True(t, err != nil, "should return error")
	})
}

func TestStartHelper_Successful(t *testing.T) {
	conHandler := mocks.ClientMockConnectionHandler{[]types.NodeAddress{":3001"}, []types.NodeAddress{}}
	t.Run("attempt start with super node online", func(t *testing.T) {
		config := Config{
			IdealNumberOfHops: 4,
			MinNumberOfHops:   2,
			KnownSuperNode:    ":3001",
		}

		c1, err := StartHelper(config, &conHandler)
		assert.True(t, err == nil, "should not return error")

		assert.Equal(t, c1, &Client{
			superNode:         ":3001",
			connHandler:       &conHandler,
			idealNumberOfHops: 4,
			minNumberOfHops:   2,
		})
	})
}

func TestGetPath(t *testing.T) {
	t.Run("get path with many nodes online", func(t *testing.T) {
		conHandler := mocks.ClientMockConnectionHandler{[]types.NodeAddress{":3001", ":3003", ":3004", ":3005", ":3006", ":3007", ":3008", ":3010", ":30011", ":30012"}, []types.NodeAddress{}}
		config := Config{
			IdealNumberOfHops: 4,
			MinNumberOfHops:   2,
			KnownSuperNode:    ":3001",
		}

		c1, _ := StartHelper(config, &conHandler)
		path := c1.getPath(config.IdealNumberOfHops)
		assert.Equal(t, config.IdealNumberOfHops, uint16(len(path)))
	})

	t.Run("get path with < IdealNumberOfHops  but > MinNumberOfHops", func(t *testing.T) {
		conHandler := mocks.ClientMockConnectionHandler{[]types.NodeAddress{":3001", ":3003", ":3004"}, []types.NodeAddress{}}
		config := Config{
			IdealNumberOfHops: 4,
			MinNumberOfHops:   2,
			KnownSuperNode:    ":3001",
		}

		c1, _ := StartHelper(config, &conHandler)
		path := c1.getPath(config.IdealNumberOfHops)
		assert.Equal(t, uint16(3), uint16(len(path)))
	})
}

func TestClient_Send(t *testing.T) {
	t.Run("get path with < IdealNumberOfHops  but > MinNumberOfHops", func(t *testing.T) {
		conHandler := mocks.ClientMockConnectionHandler{[]types.NodeAddress{":3001", ":3003", ":3004"}, []types.NodeAddress{}}
		config := Config{
			IdealNumberOfHops: 4,
			MinNumberOfHops:   2,
			KnownSuperNode:    ":3001",
		}

		c1, _ := StartHelper(config, &conHandler)

		c1.Send("hello")
		assert.Equal(t, 1, len(conHandler.SendRequests))

	})
}
