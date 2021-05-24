package mocks

import (
	"errors"
	"github.com/AdamCollins/ogre-net/internal/types"
)

// mock connection handler in which all contacted nodes respond with an error(except for self calls)
type MockConnectionHandlerAllOffline struct {
	ListenRequestLog    []types.NodeAddress
	PingRequestLog      []types.NodeAddress
	AdvertiseRequestLog []AdRequest
	NodeListenAddr      types.NodeAddress
}

func NewMockOffline(listenAddr types.NodeAddress) MockConnectionHandlerAllOffline {
	return MockConnectionHandlerAllOffline{
		ListenRequestLog:    []types.NodeAddress{},
		PingRequestLog:      []types.NodeAddress{},
		AdvertiseRequestLog: []AdRequest{},
		NodeListenAddr:      listenAddr,
	}
}

func (m *MockConnectionHandlerAllOffline) Listen(address types.NodeAddress) error {
	m.ListenRequestLog = append(m.ListenRequestLog, address)
	return nil
}

func (m *MockConnectionHandlerAllOffline) Advertise(address types.NodeAddress, knownAddresses []types.NodeAddress) error {
	ad := AdRequest{
		Target:      address,
		Advertising: knownAddresses,
	}
	m.AdvertiseRequestLog = append(m.AdvertiseRequestLog, ad)

	return nil
}

func (m *MockConnectionHandlerAllOffline) PingNode(address types.NodeAddress) error {
	// answer calls to self
	if address == m.NodeListenAddr {
		m.PingRequestLog = append(m.PingRequestLog, address)
		return nil
	}
	return errors.New("no response")
}
