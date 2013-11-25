package lex

import (
	e "errors"
)

var (
	ExpectedSpace = e.New("expected space")
	NilReader     = e.New("nil reader parameter")
)
