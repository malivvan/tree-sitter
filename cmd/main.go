package main

import (
	"context"
	sitter "github.com/malivvan/tree-sitter"
	"log"
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
	ts, err := sitter.New(ctx)
	if err != nil {
		panic(err)
	}
	p, err := ts.NewParser(ctx)
	if err != nil {
		panic(err)
	}
	// defer p.Close(ctx)

	// p.Delete()

	clang, err := ts.LanguageCpp(ctx)
	if err != nil {
		panic(err)
	}
	v, err := p.GetLanguageVersion(ctx, clang)
	if err != nil {
		panic(err)
	}
	log.Printf("c lang version: %+v\n", v)

	err = p.SetLanguage(ctx, clang)
	if err != nil {
		panic(err)
	}

	tree, err := p.ParseString(ctx, cppCode)
	if err != nil {
		panic(err)
	}
	root, err := tree.RootNode(ctx)
	if err != nil {
		panic(err)
	}
	rootKind, err := root.Kind(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node kind: %+v\n", rootKind)
	rootString, err := root.String(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node string: %+v\n", rootString)
	rootChildCount, err := root.ChildCount(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("root node child count: %+v\n", rootChildCount)

}
