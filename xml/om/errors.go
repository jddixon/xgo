package om

import (
	e "errors"
)

var (
	GraphCycleError		= e.New("graph cycle error")
	IndexOutOfBounds = e.New("index out of bounds")
	NilAttr				= e.New("nil Attr argument")
	NilHolder			= e.New("nil Holder argument")
	RuntimeError		= e.New("runtime error")
)
