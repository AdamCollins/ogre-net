package node

import "github.com/AdamCollins/ogre-net/internal/types"

type AdRequest struct {
	Target      types.NodeAddress
	Advertising []types.NodeAddress
}
type MockConnectionHandler struct {
	ListenRequestLog    []types.NodeAddress
	PingRequestLog      []types.NodeAddress
	AdvertiseRequestLog []AdRequest
}

func (m *MockConnectionHandler) Listen(address types.NodeAddress) error {
	m.ListenRequestLog = append(m.ListenRequestLog, address)
	return nil
}

func (m *MockConnectionHandler) Advertise(address types.NodeAddress, knownAddresses []types.NodeAddress) error {
	ad := AdRequest{
		Target:      address,
		Advertising: knownAddresses,
	}
	m.AdvertiseRequestLog = append(m.AdvertiseRequestLog, ad)

	return nil
}

func (m *MockConnectionHandler) PingNode(address types.NodeAddress) error {
	m.PingRequestLog = append(m.PingRequestLog, address)
	return nil
}
