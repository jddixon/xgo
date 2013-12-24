package md

import (
	e "errors"
)

var (
	EmptyHeaderTitle     = e.New("empty header title")
	EmptyID              = e.New("empty id parameter")
	EmptyURI             = e.New("empty uri parameter")
	HeaderNOutOfRange    = e.New("header N out of range")
	InvalidCharInID      = e.New("invalid char in ID")
	InvalidLineSeparator = e.New("invalid line separator (not LF, CR, zero")
	NilDocument          = e.New("nil document parameter")
	NilID                = e.New("nil id parameter")
	NilParser            = e.New("nil parser parameter")
	NilWriter            = e.New("nil writer parameter")
	NotALineSeparator    = e.New("not a line separator")
)
