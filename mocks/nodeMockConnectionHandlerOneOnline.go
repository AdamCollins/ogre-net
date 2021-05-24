package mocks

import (
	"errors"
	"github.com/AdamCollins/ogre-net/internal/types"
)

var onlineNode = ":3000"

// mock connection handler in which only node at :3000 responds
type MockConnectionHandlerOneOnline struct {
	ListenRequestLog     []types.NodeAddress
	PingResponseLog      []types.NodeAddress
	AdvertiseResponseLog []AdRequest
	NodeListenAddr       types.NodeAddress
}

func NewMockOneOnline(listenAddr types.NodeAddress) MockConnectionHandlerOneOnline {
	return MockConnectionHandlerOneOnline{
		ListenRequestLog:     []types.NodeAddress{},
		PingResponseLog:      []types.NodeAddress{},
		AdvertiseResponseLog: []AdRequest{},
		NodeListenAddr:       listenAddr,
	}
}

func (m *MockConnectionHandlerOneOnline) Listen(address types.NodeAddress) error {
	m.ListenRequestLog = append(m.ListenRequestLog, address)
	return nil
}

func (m *MockConnectionHandlerOneOnline) Advertise(address types.NodeAddress, knownAddresses []types.NodeAddress) error {

	if address != types.NodeAddress(onlineNode) && address != m.NodeListenAddr {
		return errors.New("no response")
	}
	ad := AdRequest{
		Target:        address,
		Advertisement: knownAddresses,
	}
	m.AdvertiseResponseLog = append(m.AdvertiseResponseLog, ad)

	return nil
}

func (m *MockConnectionHandlerOneOnline) PingNode(address types.NodeAddress) error {
	if address != types.NodeAddress(onlineNode) && address != m.NodeListenAddr {
		return errors.New("no response")
	}
	m.PingResponseLog = append(m.PingResponseLog, address)
	return nil
}
