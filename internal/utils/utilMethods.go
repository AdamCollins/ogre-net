package utils

import (
	"github.com/AdamCollins/ogre-net/internal/types"
	"math/rand"
	"time"
)

func ChanToSlice(c chan types.NodeAddress) []types.NodeAddress {
	s := make([]types.NodeAddress, 0)
	for i := range c {
		s = append(s, i)
	}
	return s
}

func MinUInt16(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func ShuffleNodes(nodes []types.NodeAddress) []types.NodeAddress {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
	return nodes
}

// returns true if slice contains key, otherwise false
func ContainsNodeAddress(slice []types.NodeAddress, key types.NodeAddress) bool {
	for _, v := range slice {
		if v == key {
			return true
		}
	}
	return false
}
