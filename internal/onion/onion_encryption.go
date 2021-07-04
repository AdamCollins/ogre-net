package onion

import (
	"github.com/AdamCollins/ogre-net/internal/types"
)

// take message with payload at centre. layer message for nodes outward in.
// eg. payload: 'GET', hopList: [N1,N2,N3] => N1(N2(N3('GET')))
func NewOnionMessage(payload string, hopList []types.NodeAddress) types.OnionMessage {
	baseMessage := types.OnionMessage{
		NextLayer: nil,
		Payload:   payload,
	}
	for i := len(hopList) - 1; i >= 0; i-- {
		baseMessage = Layer(baseMessage, hopList[i])
	}

	return baseMessage
}

// add a new layer to the onto the message
// eg M1 -> M2(M1)
func Layer(message types.OnionMessage, nextHop types.NodeAddress) types.OnionMessage {
	return types.OnionMessage{
		NextHop:   nextHop,
		NextLayer: &message,
	}
}

// removes a layer from the onion message
// eg. M2(M1) -> M1
func Peel(message types.OnionMessage) types.OnionMessage {
	if message.NextLayer == nil {
		return message
	}
	return *message.NextLayer
}
