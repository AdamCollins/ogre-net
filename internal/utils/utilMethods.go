package utils

import "github.com/AdamCollins/ogre-net/internal/types"

func ChanToSlice(c chan types.NodeAddress) []types.NodeAddress {
	s := make([]types.NodeAddress, 0)
	for i := range c {
		s = append(s, i)
	}
	return s
}
