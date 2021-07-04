package client

import (
	"github.com/AdamCollins/ogre-net/internal/types"
	"log"
	"net/rpc"
)

type ConnectionHandler struct {
}

func (c ConnectionHandler) Ping(node types.NodeAddress) error {
	log.Printf("[Client] pinging %v\n", node)

	conn := rpcDial(node)
	ack := false
	defer conn.Close()
	err := conn.Call("RPCReceiver.PingRPCHandler", ack, &ack)
	if err != nil {
		return err
	}

	return nil
}

func (c ConnectionHandler) AskForNodes(node types.NodeAddress, n uint16) ([]types.NodeAddress, error) {
	log.Printf("[Client] asking %v for %v nodes\n", node, n)

	conn := rpcDial(node)
	nodes := make([]types.NodeAddress, n)
	defer conn.Close()
	err := conn.Call("RPCReceiver.GetRandomNodesSubsetHandler", n, &nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (c ConnectionHandler) SendMessage(msg types.OnionMessage, firstNode types.NodeAddress) (*types.OnionMessage, error) {
	conn := rpcDial(firstNode)
	reply := types.OnionMessage{}
	defer conn.Close()
	err := conn.Call("RPCReceiver.ReceiveMessageHandler", msg, &reply)
	if err != nil {
		log.Fatalf("[Client] Could not conect to node at %v\n", firstNode)
		return nil, err
	}
	return &reply, nil
}

func rpcDial(node types.NodeAddress) *rpc.Client {
	conn, err := rpc.Dial("tcp", string(node))
	if err != nil {
		log.Printf("[Client RPC] Could not conect to node at %v\n", node)
		return nil
	}
	return conn
}
