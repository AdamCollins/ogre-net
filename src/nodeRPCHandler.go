package src

type NodeRPCHandler struct {
	node Node
}

func (handler NodeRPCHandler) Listen(address NodeAddress, callbackNode Node) error {
	handler.node = callbackNode

	// register RPC calls
	// listen on
	return nil
}

// pings node at address and returns any neighbour nodes of target not included in NodeAddress
func (NodeRPCHandler) GetNodesNeighbours(address NodeAddress, knownNodes []NodeAddress) ([]NodeAddress, error) {
	// make rpccall to address.GetNeighboursRPCHandler(knownNodes)
	return nil, nil
}

// requests ack from node to ensure that it is still online
func (handler NodeRPCHandler) PingNode(address NodeAddress) error {
	panic("implement me")
}

// handle Ping RPC Request
func (handler NodeRPCHandler) PingRPCHandler(args interface{}, results *bool) error {
	*results = true
	return nil
}

// handle Get Neighbours RPC request
func (handler NodeRPCHandler) GetNeighboursRPCHandler(args []NodeAddress, results *[]NodeAddress) error {
	*results = handler.node.GetOnlineNodesDiff(args)
	return nil
}
