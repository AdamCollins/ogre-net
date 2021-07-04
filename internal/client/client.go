package client

import (
	"errors"
	"fmt"
	"github.com/AdamCollins/ogre-net/internal/onion"
	"github.com/AdamCollins/ogre-net/internal/types"
	"github.com/AdamCollins/ogre-net/internal/utils"
	"log"
)

type Client struct {
	superNode         types.NodeAddress
	connHandler       types.IClientConnectionHandler
	idealNumberOfHops uint16 // will try to build messages with this number of hops/layers
	minNumberOfHops   uint16 // if can't acquire this number of unique hops/layers, connection will fail
}

type Config struct {
	IdealNumberOfHops uint16
	MinNumberOfHops   uint16
	KnownSuperNode    types.NodeAddress
}

func Start(config Config) (*Client, error) {
	return StartHelper(config, &ConnectionHandler{})
}

func StartHelper(config Config, connHandler types.IClientConnectionHandler) (*Client, error) {
	err := connHandler.Ping(config.KnownSuperNode)
	if err != nil {
		return nil, err
	}

	return &Client{
		superNode:         config.KnownSuperNode,
		connHandler:       connHandler,
		idealNumberOfHops: config.IdealNumberOfHops,
		minNumberOfHops:   config.MinNumberOfHops,
	}, nil

}

// queries super node for n*2 nodes.
// from these nodes are reshuffled and a sample of n is selected for the path.
func (client *Client) getPath(n uint16) []types.NodeAddress {
	// TODO this function (n*2) is chosen randomly. should determine a more thought out rule
	// get set of nodes from super node
	nodes, err := client.connHandler.AskForNodes(client.superNode, n*2)
	if err != nil {
		log.Fatal(err)
	}

	// shuffle returned nodes and select random sample of n(or less if fewer nodes are returned)
	nodes = utils.ShuffleNodes(nodes)
	maxLen := utils.MinUInt16(n, uint16(len(nodes)))
	return nodes[:maxLen]

}

// send a request to be sent through onion net. eg. "HTTP GET google.com/index.html"
func (client *Client) Send(request string) (string, error) {
	path := client.getPath(client.idealNumberOfHops)

	// check that path meets min length requirement set in config
	if uint16(len(path)) < client.minNumberOfHops {
		return "", errors.New(fmt.Sprintf("[Client]: Tried to generate path to send message but could only "+
			"find %v nodes(min required=%v). Stopping \n", len(path), client.minNumberOfHops))
	}

	// create message using nodes
	msg := onion.NewOnionMessage(request, path)

	replyPtr, _ := client.connHandler.SendMessage(msg, path[0])

	reply := *replyPtr
	for range path {
		reply = onion.Peel(reply)
	}

	return reply.Payload, nil

}
