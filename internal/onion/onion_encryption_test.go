package onion

import (
	"fmt"
	"github.com/AdamCollins/ogre-net/internal/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type OnionMessage = types.OnionMessage

type OnionMessageTest struct {
	payload  string
	hopList  []types.NodeAddress
	expected OnionMessage
}

func TestNewOnionMessage(t *testing.T) {
	tests := []OnionMessageTest{
		{
			payload: "GET",
			hopList: []types.NodeAddress{":3001", ":3002"},
			expected: OnionMessage{
				NextHop: ":3001",
				NextLayer: &OnionMessage{
					NextHop: ":3002",
					NextLayer: &OnionMessage{
						NextLayer: nil,
						Payload:   "GET",
					},
				},
			},
		},
		{
			payload: "GET",
			hopList: []types.NodeAddress{":3001", ":3002", ":3003"},
			expected: OnionMessage{
				NextHop: ":3001",
				NextLayer: &OnionMessage{
					NextHop: ":3002",
					NextLayer: &OnionMessage{
						NextHop: ":3003",
						NextLayer: &OnionMessage{
							NextLayer: nil,
							Payload:   "GET",
						},
					},
				},
			},
		},
		{
			payload: "GET",
			hopList: []types.NodeAddress{":3001"},
			expected: OnionMessage{
				NextHop: ":3001",
				NextLayer: &OnionMessage{
					NextLayer: nil,
					Payload:   "GET",
				},
			},
		},
		{
			payload: "GET",
			hopList: []types.NodeAddress{},
			expected: OnionMessage{
				NextLayer: nil,
				Payload:   "GET",
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("payload: %v. hopList: %v", test.payload, test.hopList), func(t *testing.T) {

			msg := NewOnionMessage(test.payload, test.hopList)

			if !reflect.DeepEqual(test.expected, msg) {
				t.Fatalf("%v, does not equal expected value %v", msg, test.expected)
			}
			t.Logf("Success: %v, does equal expected value %v", msg, test.expected)
		})
	}
}

func TestPeel(t *testing.T) {
	// [input, expectedOutput]
	tests := [][2]types.OnionMessage{
		{
			{
				NextHop:   ":3001",
				NextLayer: &OnionMessage{Payload: "GET", NextLayer: nil},
			},
			{Payload: "GET", NextLayer: nil},
		},
		{
			{
				NextHop: ":3002",
				NextLayer: &OnionMessage{
					NextHop:   ":3001",
					NextLayer: &OnionMessage{Payload: "GET", NextLayer: nil},
				},
			},
			{
				NextHop:   ":3001",
				NextLayer: &OnionMessage{Payload: "GET", NextLayer: nil},
			},
		},
		{
			{Payload: "GET", NextLayer: nil},
			{Payload: "GET", NextLayer: nil},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("init: %v. expected after peel: %v", test[0], test[1]), func(t *testing.T) {
			out := Peel(test[0])
			assert.Equal(t, test[1], out)
		})
	}

}
