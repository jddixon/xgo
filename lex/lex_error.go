package lex

// xgo/xml/lex/lex_error.go

import (
	"fmt"
)

type LexError struct {
	text string
}

func (lxErr *LexError) Error() string {
	return lxErr.text
}

//
func (lx *LexInput) NewLexError(msg string) (lxErr error) {
	lxMsg := fmt.Sprintf("line %d col %d: %s", lx.lineNo, lx.colNo, msg)
	lxErr = &LexError{lxMsg}
	return
}
