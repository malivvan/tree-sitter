package sitter

import (
	"fmt"
	"strings"
)

type (
	Language struct {
		t *TreeSitter
		l uint64
	}

	LanguageError struct {
		version uint64
	}
)

func (l LanguageError) Error() string {
	return fmt.Sprintf("Incompatible language version %d", l.version)
}

func NewLanguage(l uint64, t *TreeSitter) Language {
	return Language{l: l, t: t}
}

func (t *TreeSitter) Language(name string) (Language, error) {
	name = strings.ToLower(name)
	lang, ok := t.lang[name]
	if !ok {
		return Language{}, fmt.Errorf("initiating language: %s does not exist", name)
	}
	langPtr, err := lang.Call(t.ctx)
	if err != nil {
		return Language{}, fmt.Errorf("initiating language: %w", err)
	}
	return NewLanguage(langPtr[0], t), nil
}
