package sitter

import (
	"fmt"
)

type Node struct {
	t *TreeSitter
	n uint64
}

func newNode(t *TreeSitter, n uint64) Node {
	return Node{t, n}
}

func (ts *TreeSitter) allocateNode() (uint64, error) {
	// allocate tsnode 24 bytes
	nodePtr, err := ts.call(_malloc, uint64(24))
	if err != nil {
		return 0, fmt.Errorf("allocating node: %w", err)
	}
	return nodePtr[0], nil
}

func (n Node) Kind() (string, error) {
	nodeTypeStrPtr, err := n.t.call(_nodeType, n.n)
	if err != nil {
		return "", fmt.Errorf("getting node type: %w", err)
	}
	return n.t.readString(nodeTypeStrPtr[0])
}

func (n Node) Child(index uint64) (Node, error) {
	nodePtr, err := n.t.allocateNode()
	if err != nil {
		return Node{}, err
	}
	_, err = n.t.call(_nodeChild, nodePtr, n.n, index)
	if err != nil {
		return Node{}, fmt.Errorf("getting node child: %w", err)
	}
	return newNode(n.t, nodePtr), nil
}

func (n Node) NamedChild(index uint64) (Node, error) {
	nodePtr, err := n.t.allocateNode()
	if err != nil {
		return Node{}, err
	}
	_, err = n.t.call(_nodeNamedChild, nodePtr, n.n, index)
	if err != nil {
		return Node{}, fmt.Errorf("getting node child: %w", err)
	}
	return newNode(n.t, nodePtr), nil
}

func (n Node) IsError() (bool, error) {
	res, err := n.t.call(_nodeIsError, n.n)
	if err != nil {
		return false, fmt.Errorf("getting node is error: %w", err)
	}
	return res[0] == 1, nil
}

func (n Node) StartByte() (uint64, error) {
	res, err := n.t.call(_nodeStartByte, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node start byte: %w", err)
	}
	return res[0], nil
}

func (n Node) EndByte() (uint64, error) {
	res, err := n.t.call(_nodeEndByte, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node end byte: %w", err)
	}
	return res[0], nil
}

func (n Node) ChildCount() (uint64, error) {
	res, err := n.t.call(_nodeChildCount, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node child count: %w", err)
	}
	return res[0], nil
}

func (n Node) NamedChildCount() (uint64, error) {
	res, err := n.t.call(_nodeNamedChildCount, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node child count: %w", err)
	}
	return res[0], nil
}

func (n Node) String() (string, error) {
	strPtr, err := n.t.call(_nodeString, n.n)
	if err != nil {
		return "", fmt.Errorf("getting node string: %w", err)
	}
	return n.t.readString(strPtr[0])
}
