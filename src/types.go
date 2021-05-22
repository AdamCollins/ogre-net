package src

type NodeAddress string

type INodeConnectionHandler interface {
	Listen(address NodeAddress, callbackNode *Node) error
	Advertise(address NodeAddress, knownAddresses []NodeAddress) error
	PingNode(address NodeAddress) error
}
