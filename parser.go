package sitter

import (
	"fmt"
)

type Parser struct {
	t *TreeSitter
	p uint64
}

func (ts *TreeSitter) NewParser() (*Parser, error) {
	p, err := ts.call(_parserNew)
	if err != nil {
		return nil, fmt.Errorf("creating parser: %w", err)
	}

	return &Parser{
		t: ts,
		p: p[0],
	}, nil
}

func (p *Parser) Close() error {
	_, err := p.t.call(_parserDelete, p.p)
	if err != nil {
		return fmt.Errorf("closing parser: %w", err)
	}
	return nil
}

func (p *Parser) SetLanguage(l *Language) error {
	ok, err := p.t.call(_parserSetLanguage, p.p, l.l)
	if err != nil {
		return fmt.Errorf("setting language: %w", err)
	}
	if ok[0] == 0 {
		v, err := p.GetLanguageVersion(l)
		if err != nil {
			return err
		}
		return LanguageError{v}
	}
	return nil
}

func (p *Parser) GetLanguageVersion(l *Language) (uint64, error) {
	v, err := p.t.call(_languageVersion, l.l)
	if err != nil {
		return 0, fmt.Errorf("getting language version: %w", err)
	}
	return v[0], nil
}

func (p *Parser) ParseString(str string) (*Tree, error) {
	strPtr, strSize, freeStr, err := p.t.allocateString(str)
	defer freeStr()

	tree, err := p.t.call(_parserParseString, p.p, uint64(0), strPtr, strSize)
	if err != nil {
		return nil, fmt.Errorf("calling ts_parser_parse_string: %w", err)
	}
	return newTree(p.t, tree[0]), nil
}
