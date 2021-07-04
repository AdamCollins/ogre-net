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

func (connectionHandler *ConnectionHandler) SetCallbackNode(node *Node) {
	connectionHandler.receiver.node = node
}

func (connectionHandler *ConnectionHandler) Listen(address types.NodeAddress) error {
	if connectionHandler.receiver.node == nil {
		return errors.New("cannot listen to conHandler before SetCallbackNode(node *Node) is called")
	}
	server := rpc.NewServer()
	err := server.Register(connectionHandler.receiver)

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
func (connectionHandler ConnectionHandler) PingNode(address types.NodeAddress) error {
	// don't bother taking up network time if targeting self
	if connectionHandler.receiver.node.ListenAddr == address {
		return nil
	}

	log.Printf("[Node %v] pinging %v\n", connectionHandler.receiver.node.NodeId, address)
	conn, err := rpc.Dial("tcp", string(address))
	if err != nil {
		log.Printf("[Node RPC] %v Could not conect to node at %v", connectionHandler.receiver.node.NodeId, address)
		return err
	}
	defer conn.Close()
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
func (connectionHandler ConnectionHandler) Advertise(target types.NodeAddress, advertising []types.NodeAddress) error {
	// make rpccall to target.AdvertiseHandler(advertising)
	conn, err := rpc.Dial("tcp", string(target))
	if err != nil {
		log.Printf("[Node RPC] %v Could not conect to node at %v", connectionHandler.receiver.node.NodeId, target)
		return err
	}
	defer conn.Close()
	ack := false
	err = conn.Call("RPCReceiver.AdvertiseHandler", advertising, &ack)
	if err != nil {
		return err
	}

	return nil
}

// handle Get Neighbours RPC request
func (receiver RPCReceiver) AdvertiseHandler(args []types.NodeAddress, results *bool) error {

	go receiver.node.ReceiveAdvertisement(args)

	*results = true
	return nil
}

// handle calls to get Nodes from client
//RPCReceiver.GetRandomNodesSubsetHandler
func (receiver RPCReceiver) GetRandomNodesSubsetHandler(args uint16, results *[]types.NodeAddress) error {
	nodes := receiver.node.GetRandomNodesSubset(args)
	*results = nodes
	return nil
}

func (receiver RPCReceiver) ReceiveMessageHandler(args types.OnionMessage, results *types.OnionMessage) error {
	*results = receiver.node.ReceiveMessage(args)
	return nil
}

// takes OnionMessage, msg, and forwards it to address
// returns the response from the node at address in the form of an OnionMessage
func (connectionHandler ConnectionHandler) ForwardMessage(msg types.OnionMessage, address types.NodeAddress) (types.OnionMessage, error) {
	conn, err := rpc.Dial("tcp", string(address))
	if err != nil {
		log.Fatalf("[Node RPC] %v Could not conect to node at %v", connectionHandler.receiver.node.NodeId, address)
		return types.OnionMessage{}, err
	}
	reply := types.OnionMessage{}
	defer conn.Close()
	err = conn.Call("RPCReceiver.ReceiveMessageHandler", msg, &reply)
	if err != nil {
		log.Fatalf("[Node RPC] %v Could not conect to node at %v", connectionHandler.receiver.node.NodeId, address)
		return types.OnionMessage{}, err
	}

	return reply, nil
}
