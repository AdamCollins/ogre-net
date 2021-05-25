package node_set

import (
	"fmt"
	"github.com/AdamCollins/ogre-net/internal/types"
	"reflect"
	"testing"
)

type Test struct {
	input    []types.NodeAddress
	expected []types.NodeAddress
}

type DiffTest struct {
	set      []types.NodeAddress
	diff     []types.NodeAddress
	expected []types.NodeAddress
}

func TestNodeSet_NewNodeSet(t *testing.T) {
	nodeset := NewNodeSet()
	entries := nodeset.GetNodes()
	if len(entries) != 0 {
		t.Fatalf("set was not initialized correctly")
	}
}

func TestNodeSet_AddOnlineNodes(t *testing.T) {
	tests := []Test{
		{[]types.NodeAddress{":3001"}, []types.NodeAddress{":3001"}},
		{[]types.NodeAddress{":3001", ":3002"}, []types.NodeAddress{":3001", ":3002"}},
		{[]types.NodeAddress{}, []types.NodeAddress{}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			nodeset := NewNodeSet()
			nodeset.AddOnlineNodes(test.input)
			result := nodeset.GetNodes()
			if len(result) != len(test.expected) {
				t.Fatalf("%s, does not equal expected value %s", result, test.expected)
			}
		})
	}
}

func TestNodeSet_AddOnlineNodesOneAtATime(t *testing.T) {
	tests := []Test{
		{[]types.NodeAddress{":3001"}, []types.NodeAddress{":3001"}},
		{[]types.NodeAddress{":3001", ":3002"}, []types.NodeAddress{":3001", ":3002"}},
		{[]types.NodeAddress{}, []types.NodeAddress{}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			nodeset := NewNodeSet()
			for _, addr := range test.input {
				nodeset.AddOnlineNode(addr)
			}
			result := nodeset.GetNodes()
			if len(result) != len(test.expected) {
				t.Fatalf("%s, does not equal expected value %s", result, test.expected)
			}
		})
	}
}

func TestNodeSet_NotInSet(t *testing.T) {

	// eg. diffNodes = [node1, node 2, node3], set=[node 1, node 2]

	tests := []DiffTest{
		{[]types.NodeAddress{":3001"}, []types.NodeAddress{":3001"}, []types.NodeAddress{}},
		{[]types.NodeAddress{":3001", ":3002"}, []types.NodeAddress{":3001", ":3002", ":3003"}, []types.NodeAddress{":3003"}},
		{[]types.NodeAddress{":3001", ":3002", ":3003"}, []types.NodeAddress{":3001", ":3002"}, []types.NodeAddress{}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("set:%v - diff:%v, should return:%v", test.set, test.diff, test.expected), func(t *testing.T) {
			nodeset := NewNodeSet()
			nodeset.AddOnlineNodes(test.set)
			diff := nodeset.GetDifference(test.diff)

			if !reflect.DeepEqual(test.expected, diff) {
				t.Fatalf("%s, does not equal expected value %s", diff, test.expected)
			}
		})
	}
}

type RemoveTest struct {
	Initial  []types.NodeAddress
	Remove   []types.NodeAddress
	Expected []types.NodeAddress
}

func TestNodeSet_RemoveOnlineNodes(t *testing.T) {
	tests := []RemoveTest{
		{
			[]types.NodeAddress{},
			[]types.NodeAddress{},
			[]types.NodeAddress{},
		},
		{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{},
			[]types.NodeAddress{":3001"},
		},
		{
			[]types.NodeAddress{":3001", ":3002"},
			[]types.NodeAddress{":3002"},
			[]types.NodeAddress{":3001"},
		},
		{
			[]types.NodeAddress{":3001"},
			[]types.NodeAddress{":3003"},
			[]types.NodeAddress{":3001"},
		},
		{
			[]types.NodeAddress{},
			[]types.NodeAddress{":3003"},
			[]types.NodeAddress{},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("set:(%v) removing(%v), should return:%v", test.Initial, test.Remove, test.Expected), func(t *testing.T) {
			nodeset := NewNodeSet()
			nodeset.AddOnlineNodes(test.Initial)
			nodeset.RemoveOnlineNodes(test.Remove)

			diff := nodeset.GetNodes()
			if !reflect.DeepEqual(test.Expected, diff) {
				t.Fatalf("%s, does not equal expected value %s", diff, test.Expected)
			}
		})
	}
}
