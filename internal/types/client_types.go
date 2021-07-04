package types

type IClientConnectionHandler interface {
	Ping(node NodeAddress) error                                                // pings 'node'
	AskForNodes(node NodeAddress, n uint16) ([]NodeAddress, error)              // asks 'node' for up to n nodes
	SendMessage(msg OnionMessage, firstNode NodeAddress) (*OnionMessage, error) // sends an Onion encoded message to 'firstNode'
}
