package node

import (
	"errors"
	"github.com/AdamCollins/ogre-net/internal/types"
	"log"
	"net"
	"net/rpc"
)

type ConnectionHandler struct {
	receiver RPCReceiver
}

type RPCReceiver struct {
	node *Node
}

func (manager *ConnectionHandler) SetCallbackNode(node *Node) {
	manager.receiver.node = node
}

func (manager *ConnectionHandler) Listen(address types.NodeAddress) error {
	if manager.receiver.node == nil {
		return errors.New("cannot listen to conHandler before SetCallbackNode(node *Node) is called")
	}
	server := rpc.NewServer()
	err := server.Register(manager.receiver)

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
func (manager ConnectionHandler) PingNode(address types.NodeAddress) error {
	log.Printf("[Node %v] pinging %v\n", manager.receiver.node.NodeId, address)
	conn, err := rpc.Dial("tcp", string(address))
	if err != nil {
		log.Printf("[Node RPC] %v Could not conect to node at %v", manager.receiver.node.NodeId, address)
		return err
	}
	ack := false
	err = conn.Call("RPCReceiver.PingRPCHandler", ack, &ack)
	if err != nil {
		return err
	}

	return nil
}

// handle Ping RPC Request
func (receiver RPCReceiver) PingRPCHandler(args bool, results *bool) error {
	*results = true
	log.Printf("[Node %v] pong!\n", receiver.node.NodeId)

	return nil
}

// pings node at address and returns any neighbour nodes of target not included in NodeAddress
func (manager ConnectionHandler) Advertise(target types.NodeAddress, advertising []types.NodeAddress) error {
	// make rpccall to target.AdvertiseHandler(advertising)
	conn, err := rpc.Dial("tcp", string(target))
	if err != nil {
		log.Printf("[Node RPC] %v Could not conect to node at %v", manager.receiver.node.NodeId, target)
		return err
	}
	ack := false
	err = conn.Call("RPCReceiver.AdvertiseHandler", advertising, &ack)
	if err != nil {
		return err
	}

	return nil
}

// handle Get Neighbours RPC request
func (receiver RPCReceiver) AdvertiseHandler(args []types.NodeAddress, results *bool) error {

	go receiver.node.CheckForNewNodes(args)

	*results = true
	return nil
}
