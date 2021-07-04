package types

type NodeAddress string

type PublicKey string
type PrivateKey string

type INodeConnectionHandler interface {
	Listen(address NodeAddress) error
	Advertise(address NodeAddress, knownAddresses []NodeAddress) error
	PingNode(address NodeAddress) error
	ForwardMessage(msg OnionMessage, address NodeAddress) (OnionMessage, error)
}
