package src

type NodeAddress string

type INodeRPCHandler interface {
	Listen(address NodeAddress, callbackNode Node) error
	GetNodesNeighbours(address NodeAddress, knownAddresses []NodeAddress) ([]NodeAddress, error)
	PingNode(address NodeAddress) error
}
