package main

import (
	"context"
	"fmt"
	sitter "github.com/malivvan/tree-sitter"
	"log"
	"strings"
)

var cppCode = `#include <iostream>

int main() {
    std::cout << "Hello World!";
    return 0;
}`

var cCode = `#include <stdio.h>

int main() {
	printf("Hello World!");
	return 0;
}`

func main() {
	ctx := context.Background()
	ts, err := sitter.New(nil, nil)
	if err != nil {
		panic(err)
	}
	p, err := ts.NewParser()
	if err != nil {
		panic(err)
	}
	// defer p.Close(ctx)

	// p.Delete()

	clang, err := ts.Language("cpp")
	if err != nil {
		panic(err)
	}
	v, err := p.GetLanguageVersion(clang)
	if err != nil {
		panic(err)
	}
	log.Printf("c lang version: %+v\n", v)

	err = p.SetLanguage(clang)
	if err != nil {
		panic(err)
	}
	tree, err := p.ParseString(cppCode)
	if err != nil {
		panic(err)
	}
	root, err := tree.RootNode()
	if err != nil {
		panic(err)
	}
	err = walk(ctx, root, 0)
	if err != nil {
		panic(err)
	}
}

func walk(ctx context.Context, node sitter.Node, indent int) error {
	s, err := node.String()
	if err != nil {
		return err
	}
	fmt.Println(strings.Repeat("  ", indent) + s)
	n, err := node.ChildCount()
	if err != nil {
		return err
	}
	for i := uint64(0); i < n; i++ {
		child, err := node.Child(i)
		if err != nil {
			return err
		}
		err = walk(ctx, child, indent+1)
		if err != nil {
			return err
		}
	}
	return nil
}
