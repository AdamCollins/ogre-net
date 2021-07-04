package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinUInt16(t *testing.T) {
	t.Run("min(1,2)", func(t *testing.T) {
		val := MinUInt16(1, 2)
		assert.Equal(t, uint16(1), val)
	})

	t.Run("min(5,3)", func(t *testing.T) {
		val := MinUInt16(5, 3)
		assert.Equal(t, uint16(3), val)
	})

	t.Run("min(2,1)", func(t *testing.T) {
		val := MinUInt16(2, 1)
		assert.Equal(t, uint16(1), val)
	})

	t.Run("min(10,10)", func(t *testing.T) {
		val := MinUInt16(10, 10)
		assert.Equal(t, uint16(10), val)
	})

	t.Run("min(65535,1)", func(t *testing.T) {
		val := MinUInt16(65535, 1)
		assert.Equal(t, uint16(1), val)
	})
}

func TestShuffleNodes(t *testing.T) {

}
