package onion

import (
	"fmt"
	"github.com/AdamCollins/ogre-net/internal/types"
	"reflect"
	"testing"
)

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
