package mocks

import "github.com/AdamCollins/ogre-net/internal/types"

type AdRequest struct {
	Target      types.NodeAddress
	Advertising []types.NodeAddress
}

// mock connection handler in which all requested contacted nodes respond to RPC calls
type MockConnectionHandlerAllOnline struct {
	ListenRequestLog    []types.NodeAddress
	PingRequestLog      []types.NodeAddress
	AdvertiseRequestLog []AdRequest
}

func (m *MockConnectionHandlerAllOnline) Listen(address types.NodeAddress) error {
	m.ListenRequestLog = append(m.ListenRequestLog, address)
	return nil
}

func (m *MockConnectionHandlerAllOnline) Advertise(address types.NodeAddress, knownAddresses []types.NodeAddress) error {
	ad := AdRequest{
		Target:      address,
		Advertising: knownAddresses,
	}
	m.AdvertiseRequestLog = append(m.AdvertiseRequestLog, ad)

	return nil
}

func (m *MockConnectionHandlerAllOnline) PingNode(address types.NodeAddress) error {
	m.PingRequestLog = append(m.PingRequestLog, address)
	return nil
}
