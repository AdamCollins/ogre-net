package main

import (
	"flag"
	"github.com/AdamCollins/ogre-net/internal/node"
	"github.com/AdamCollins/ogre-net/internal/types"
	"log"
	"strings"
)

func main() {

	// get Node Id
	var config node.Config
	var knownNodeStr string
	// get known Nodes
	var listenAddrStr string

	flag.StringVar(&config.NodeId, "id", config.NodeId, "Node ID, e.g. node1")
	flag.StringVar(&knownNodeStr, "known", knownNodeStr, "Known Nodes, e.g. :3001 or :3001,:3002")
	flag.StringVar(&listenAddrStr, "listen", listenAddrStr, "Listen Address, e.g. :3000")
	flag.Parse()
	if len(listenAddrStr) == 0 || len(knownNodeStr) == 0 || len(config.NodeId) == 0 {
		log.Fatal("Could not parse params")
		return
	}

	// get known nodes
	var knownNodes []types.NodeAddress
	for _, addr := range strings.Split(knownNodeStr, ",") {
		knownNodes = append(knownNodes, types.NodeAddress(addr))
	}
	config.KnownNodes = knownNodes

	// get the address to listen on
	config.ListenAddr = types.NodeAddress(listenAddrStr)

	node.StartListening(config)

}
