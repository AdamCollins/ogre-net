package types

type NodeAddress string

type INodeConnectionHandler interface {
	Listen(address NodeAddress) error
	Advertise(address NodeAddress, knownAddresses []NodeAddress) error
	PingNode(address NodeAddress) error
}
