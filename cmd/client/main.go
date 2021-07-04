package main

import (
	"github.com/AdamCollins/ogre-net/internal/client"
	"log"
)

func main() {

	config := client.Config{
		IdealNumberOfHops: 4,
		MinNumberOfHops:   2,
		KnownSuperNode:    ":3001",
	}

	c1, err := client.Start(config)
	if err != nil {
		log.Fatal(err)
	}

	reply, err := c1.Send("GET google.ca/index.html")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("[Client] reply: %v\n", reply)

}
