package sitter

import (
	"fmt"
)

type Tree struct {
	ts *TreeSitter
	t  uint64
}

func newTree(ts *TreeSitter, t uint64) *Tree {
	return &Tree{ts, t}
}

func (t *Tree) RootNode() (*Node, error) {
	// allocate tsnode 24 bytes
	nodePtr, err := t.ts.call(_malloc, 24)
	if err != nil {
		return nil, fmt.Errorf("allocating node: %w", err)
	}

	_, err = t.ts.call(_treeRootNode, nodePtr[0], t.t)
	if err != nil {
		return nil, fmt.Errorf("getting tree root node: %w", err)
	}
	return newNode(t.ts, nodePtr[0]), nil
}
