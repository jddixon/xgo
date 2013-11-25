package template

import (
	e "errors"
)

var (
	NilContext = e.New("nil context argument")
	NilReader  = e.New("nil reader argument")
	NilWriter  = e.New("nil writer argument")
)
