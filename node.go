package sitter

import (
	"fmt"
)

type Node struct {
	ts  *TreeSitter
	ptr uint64
}

func newNode(ts *TreeSitter, n uint64) *Node {
	return &Node{ts, n}
}

func (ts *TreeSitter) allocateNode() (uint64, error) {
	// allocate tsnode 24 bytes
	nodePtr, err := ts.call(_malloc, uint64(24))
	if err != nil {
		return 0, fmt.Errorf("allocating node: %w", err)
	}
	return nodePtr[0], nil
}

func (n *Node) Kind() (string, error) {
	nodeTypeStrPtr, err := n.ts.call(_nodeType, n.ptr)
	if err != nil {
		return "", fmt.Errorf("getting node type: %w", err)
	}
	return n.ts.readString(nodeTypeStrPtr[0])
}

func (n *Node) Child(index uint64) (*Node, error) {
	nodePtr, err := n.ts.allocateNode()
	if err != nil {
		return nil, err
	}
	_, err = n.ts.call(_nodeChild, nodePtr, n.ptr, index)
	if err != nil {
		return nil, fmt.Errorf("getting node child: %w", err)
	}
	return newNode(n.ts, nodePtr), nil
}

func (n *Node) NamedChild(index uint64) (*Node, error) {
	nodePtr, err := n.ts.allocateNode()
	if err != nil {
		return nil, err
	}
	_, err = n.ts.call(_nodeNamedChild, nodePtr, n.ptr, index)
	if err != nil {
		return nil, fmt.Errorf("getting node child: %w", err)
	}
	return newNode(n.ts, nodePtr), nil
}

func (n *Node) IsError() (bool, error) {
	res, err := n.ts.call(_nodeIsError, n.ptr)
	if err != nil {
		return false, fmt.Errorf("getting node is error: %w", err)
	}
	return res[0] == 1, nil
}

func (n *Node) StartByte() (uint64, error) {
	res, err := n.ts.call(_nodeStartByte, n.ptr)
	if err != nil {
		return 0, fmt.Errorf("getting node start byte: %w", err)
	}
	return res[0], nil
}

func (n *Node) EndByte() (uint64, error) {
	res, err := n.ts.call(_nodeEndByte, n.ptr)
	if err != nil {
		return 0, fmt.Errorf("getting node end byte: %w", err)
	}
	return res[0], nil
}

func (n *Node) ChildCount() (uint64, error) {
	res, err := n.ts.call(_nodeChildCount, n.ptr)
	if err != nil {
		return 0, fmt.Errorf("getting node child count: %w", err)
	}
	return res[0], nil
}

func (n *Node) NamedChildCount() (uint64, error) {
	res, err := n.ts.call(_nodeNamedChildCount, n.ptr)
	if err != nil {
		return 0, fmt.Errorf("getting node child count: %w", err)
	}
	return res[0], nil
}

func (n *Node) String() (string, error) {
	strPtr, err := n.ts.call(_nodeString, n.ptr)
	if err != nil {
		return "", fmt.Errorf("getting node string: %w", err)
	}
	return n.ts.readString(strPtr[0])
}
