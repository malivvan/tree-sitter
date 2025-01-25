package sitter

import (
	"context"
	"fmt"
)

type (
	Language struct {
		t TreeSitter
		l uint64
	}

	LanguageError struct {
		version uint64
	}
)

func (l LanguageError) Error() string {
	return fmt.Sprintf("Incompatible language version %d", l.version)
}

func NewLanguage(l uint64, t TreeSitter) Language {
	return Language{l: l, t: t}
}

func (t TreeSitter) LanguageC(ctx context.Context) (Language, error) {
	cLangPtr, err := t.languageC.Call(ctx)
	if err != nil {
		return Language{}, fmt.Errorf("initiating c language: %w", err)
	}
	return NewLanguage(cLangPtr[0], t), nil
}

func (t TreeSitter) LanguageCpp(ctx context.Context) (Language, error) {
	cLangPtr, err := t.languageCpp.Call(ctx)
	if err != nil {
		return Language{}, fmt.Errorf("initiating cpp language: %w", err)
	}
	return NewLanguage(cLangPtr[0], t), nil
}
