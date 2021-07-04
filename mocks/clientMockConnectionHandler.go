package mocks

import (
	"errors"
	"github.com/AdamCollins/ogre-net/internal/types"
)

type ClientMockConnectionHandler struct {
	OnlineNodes  []types.NodeAddress
	SendRequests []types.NodeAddress
}

func (m ClientMockConnectionHandler) Ping(node types.NodeAddress) error {
	for _, e := range m.OnlineNodes {
		if e == node {
			return nil
		}
	}

	return errors.New("could not reach " + string(node))
}

func (m ClientMockConnectionHandler) AskForNodes(node types.NodeAddress, n uint16) ([]types.NodeAddress, error) {
	return m.OnlineNodes, nil
}

func (m *ClientMockConnectionHandler) SendMessage(msg types.OnionMessage, firstNode types.NodeAddress) (*types.OnionMessage, error) {
	m.SendRequests = append(m.SendRequests, firstNode)
	return &msg, nil
}
