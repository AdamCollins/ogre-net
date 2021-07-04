package mocks

import (
	"errors"
	"github.com/AdamCollins/ogre-net/internal/types"
	"github.com/AdamCollins/ogre-net/internal/utils"
)

type INode interface {
	ReceiveMessage(msg types.OnionMessage) types.OnionMessage
}

type AdRequest struct {
	Target        types.NodeAddress
	Advertisement []types.NodeAddress
}

// mock connection handler in which only node at :3000 responds
type NodeMockConnectionHandler struct {
	ListenRequestLog     []types.NodeAddress
	PingResponseLog      []types.NodeAddress
	AdvertiseResponseLog []AdRequest
	ForwardRequestLog    []types.NodeAddress
	NodeListenAddr       types.NodeAddress
	OnlineNodes          []types.NodeAddress
	Node                 INode
}

func (m *NodeMockConnectionHandler) ForwardMessage(msg types.OnionMessage, address types.NodeAddress) (types.OnionMessage, error) {
	m.ForwardRequestLog = append(m.ForwardRequestLog, address)
	res := m.Node.ReceiveMessage(msg)
	res.NextHop = address // inject address to simulate multiple nodes
	return res, nil
}

func (m *NodeMockConnectionHandler) Listen(address types.NodeAddress) error {
	m.ListenRequestLog = append(m.ListenRequestLog, address)
	return nil
}

func (m *NodeMockConnectionHandler) Advertise(address types.NodeAddress, knownAddresses []types.NodeAddress) error {

	// return error if node is not online and not self call
	if !utils.ContainsNodeAddress(m.OnlineNodes, address) && address != m.NodeListenAddr {
		return errors.New("no response")
	}

	ad := AdRequest{
		Target:        address,
		Advertisement: knownAddresses,
	}
	m.AdvertiseResponseLog = append(m.AdvertiseResponseLog, ad)

	return nil
}

func (m *NodeMockConnectionHandler) PingNode(address types.NodeAddress) error {

	// return error if node is not online and not self call
	if !utils.ContainsNodeAddress(m.OnlineNodes, address) && address != m.NodeListenAddr {
		return errors.New("no response")
	}
	m.PingResponseLog = append(m.PingResponseLog, address)
	return nil
}
