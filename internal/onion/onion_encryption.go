package onion

import (
	"github.com/AdamCollins/ogre-net/internal/types"
)

type OnionMessage struct {
	NextHop   types.NodeAddress // address of nextNode
	NextLayer *OnionMessage     // NextLayer down in package
	Payload   string            // data payload. only to be used on bottom layer
}

// take message with payload at centre. layer message for nodes outward in.
// eg. payload: 'GET', hopList: [N1,N2,N3] => N1(N2(N3('GET')))
func NewOnionMessage(payload string, hopList []types.NodeAddress) OnionMessage {
	baseMessage := OnionMessage{
		NextLayer: nil,
		Payload:   payload,
	}
	for i := len(hopList) - 1; i >= 0; i-- {
		baseMessage = layer(baseMessage, hopList[i])
	}

	return baseMessage
}

// add a new layer to the onto the message
// eg M1 -> M2(M1)
func layer(message OnionMessage, nextHop types.NodeAddress) OnionMessage {
	return OnionMessage{
		NextHop:   nextHop,
		NextLayer: &message,
	}
}

// removes a layer from the onion message
// eg. M2(M1) -> M1
func peel(message OnionMessage) OnionMessage {
	return *message.NextLayer
}
