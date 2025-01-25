package sitter

import (
	"context"
	"fmt"
)

type Node struct {
	t TreeSitter
	n uint64
}

func newNode(t TreeSitter, n uint64) Node {
	return Node{t, n}
}

func (t TreeSitter) allocateNode(ctx context.Context) (uint64, error) {
	// allocate tsnode 24 bytes
	nodePtr, err := t.malloc.Call(ctx, uint64(24))
	if err != nil {
		return 0, fmt.Errorf("allocating node: %w", err)
	}
	return nodePtr[0], nil
}

func (n Node) Kind(ctx context.Context) (string, error) {
	nodeTypeStrPtr, err := n.t.nodeType.Call(ctx, n.n)
	if err != nil {
		return "", fmt.Errorf("getting node type: %w", err)
	}
	return n.t.readString(ctx, nodeTypeStrPtr[0])
}

func (n Node) Child(ctx context.Context, index uint64) (Node, error) {
	nodePtr, err := n.t.allocateNode(ctx)
	if err != nil {
		return Node{}, err
	}
	_, err = n.t.nodeChild.Call(ctx, nodePtr, n.n, index)
	if err != nil {
		return Node{}, fmt.Errorf("getting node child: %w", err)
	}
	return newNode(n.t, nodePtr), nil
}

func (n Node) NamedChild(ctx context.Context, index uint64) (Node, error) {
	nodePtr, err := n.t.allocateNode(ctx)
	if err != nil {
		return Node{}, err
	}
	_, err = n.t.nodeNamedChild.Call(ctx, nodePtr, n.n, index)
	if err != nil {
		return Node{}, fmt.Errorf("getting node child: %w", err)
	}
	return newNode(n.t, nodePtr), nil
}

func (n Node) IsError(ctx context.Context) (bool, error) {
	res, err := n.t.nodeIsError.Call(ctx, n.n)
	if err != nil {
		return false, fmt.Errorf("getting node is error: %w", err)
	}
	return res[0] == 1, nil
}

func (n Node) StartByte(ctx context.Context) (uint64, error) {
	res, err := n.t.nodeStartByte.Call(ctx, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node start byte: %w", err)
	}
	return res[0], nil
}

func (n Node) EndByte(ctx context.Context) (uint64, error) {
	res, err := n.t.nodeEndByte.Call(ctx, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node end byte: %w", err)
	}
	return res[0], nil
}

func (n Node) ChildCount(ctx context.Context) (uint64, error) {
	res, err := n.t.nodeChildCount.Call(ctx, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node child count: %w", err)
	}
	return res[0], nil
}

func (n Node) NamedChildCount(ctx context.Context) (uint64, error) {
	res, err := n.t.nodeChildCount.Call(ctx, n.n)
	if err != nil {
		return 0, fmt.Errorf("getting node child count: %w", err)
	}
	return res[0], nil
}

func (n Node) String(ctx context.Context) (string, error) {
	strPtr, err := n.t.nodeString.Call(ctx, n.n)
	if err != nil {
		return "", fmt.Errorf("getting node string: %w", err)
	}
	return n.t.readString(ctx, strPtr[0])
}
