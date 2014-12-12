// Package trie is an implementation of a trie (prefix tree) data structure mapping []bytes to ints. It
// provides a small and simple API for usage as a set as well as a 'Node' API for walking the trie.
package trie

import "math"

type NodeId uint32

const (
	NilNode  NodeId = 0
	RootNode NodeId = 1
)

var (
	maxNodeId NodeId = math.MaxUint32 - 32
)

// A Node represents a logical vertex in the trie structure.
type node struct {
	branches [256]NodeId
	val      int
	terminal bool
}

// A Trie is a a prefix tree.
type Trie struct {
	nodes []node
}

// New construct a new, empty Trie ready for use.
func New() *Trie {
	return &Trie{
		// reserve the first node so that node ids start with 1 so that we can
		// use 0 to mean 'nil' and allocate root node
		nodes: make([]node, 2, 512),
	}
}

// Put inserts the mapping k -> v into the Trie, overwriting any previous value.
// It returns true if the element was not previously in t.
func (t *Trie) Put(k []byte, v int) bool {
	nodeId := RootNode
	for _, c := range k {
		nextNodeId := t.Walk(nodeId, c)
		if nextNodeId == NilNode {
			nextNodeId = NodeId(len(t.nodes))
			if nextNodeId > maxNodeId {
				panic("too many nodes")
			}
			t.nodes = append(t.nodes, node{})
			currNode := &t.nodes[int(nodeId)]
			currNode.branches[c] = nextNodeId
		}
		nodeId = nextNodeId
	}
	node := &t.nodes[int(nodeId)]
	node.val = v
	if node.terminal {
		return false
	}
	node.terminal = true
	return true
}

// Get the value corresponding to k in t, if any.
func (t *Trie) Get(k []byte) (v int, ok bool) {
	nodeId := RootNode
	for _, c := range k {
		next := t.Walk(nodeId, c)
		if next == NilNode {
			return 0, false
		}
		nodeId = next
	}
	node := &t.nodes[int(nodeId)]
	if node.terminal {
		return node.val, true
	}
	return 0, false
}

// Walk returns the node reached along edge c, if one exists. If node doesn't
// exist we return nil
func (t *Trie) Walk(nodeId NodeId, c byte) NodeId {
	node := &t.nodes[int(nodeId)]
	return node.branches[int(c)]
}

// Terminal indicates whether n is terminal in the trie (that is, whether the path from the root to n
// represents an element in the set). For instance, if the root node is terminal, then []byte{} is in the
// trie.
func (t *Trie) Terminal(nodeId NodeId) bool {
	node := &t.nodes[int(nodeId)]
	return node.terminal
}

// Val gives the value associated with this node. It panics if n is not terminal.
func (t *Trie) Val(nodeId NodeId) int {
	node := &t.nodes[int(nodeId)]
	if !node.terminal {
		panic("Val called on non-terminal node")
	}
	return node.val
}
