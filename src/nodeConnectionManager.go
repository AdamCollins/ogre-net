package src

import (
	"log"
	"net"
	"net/rpc"
)

type NodeConnectionManager struct {
	node *Node
}

func (handler *NodeConnectionManager) Listen(address NodeAddress, callbackNode *Node) error {
	handler.node = callbackNode
	server := rpc.NewServer()
	err := server.Register(handler)

	if err != nil {
		return err
	}

	connectionListener, err := net.Listen("tcp", string(address))

	if err != nil {
		return err
	}

	go server.Accept(connectionListener)

	return nil
}

// requests ack from node to ensure that it is still online
func (handler NodeConnectionManager) PingNode(address NodeAddress) error {
	log.Printf("[Node %v] pinging %v\n", handler.node.NodeId, address)
	conn, err := rpc.Dial("tcp", string(address))
	if err != nil {
		log.Printf("[Node RPC] %v Could not conect to node at %v", handler.node.NodeId, address)
		return err
	}
	ack := false
	err = conn.Call("NodeConnectionManager.PingRPCHandler", ack, &ack)
	if err != nil {
		return err
	}

	return nil
}

// handle Ping RPC Request
func (handler NodeConnectionManager) PingRPCHandler(args bool, results *bool) error {
	*results = true
	log.Printf("[Node %v] pong!\n", handler.node.NodeId)

	return nil
}

// pings node at address and returns any neighbour nodes of target not included in NodeAddress
func (handler NodeConnectionManager) Advertise(address NodeAddress, knownNodes []NodeAddress) error {
	// make rpccall to address.AdvertiseRPCHandler(knownNodes)
	conn, err := rpc.Dial("tcp", string(address))
	if err != nil {
		log.Printf("[Node RPC] %v Could not conect to node at %v", handler.node.NodeId, address)
		return err
	}
	ack := false
	err = conn.Call("NodeConnectionManager.AdvertiseRPCHandler", knownNodes, &ack)
	if err != nil {
		return err
	}

	return nil
}

// handle Get Neighbours RPC request
func (handler NodeConnectionManager) AdvertiseRPCHandler(args []NodeAddress, results *bool) error {

	go handler.node.CheckForNewNodes(args)

	*results = true
	return nil
}
