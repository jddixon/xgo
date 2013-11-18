package om

import (
	e "errors"
)

var (
	EmptyNamespace   = e.New("empty namespace argument")
	EmptyOtherList   = e.New("other list is empty")
	GraphCycleError  = e.New("graph cycle error")
	IllFormedPI      = e.New("ill-formed processing instruction")
	IndexOutOfBounds = e.New("index out of bounds")
	NilAttr          = e.New("nil Attr argument")
	NilChild         = e.New("nil child argument")
	NilDocType       = e.New("nil DocType argument")
	NilDocument      = e.New("nil document argument")
	NilHolder        = e.New("nil Holder argument")
	NilNode          = e.New("nil Node argument")
	NodeListNotEmpty = e.New("NodeList is not empty")
	RuntimeError     = e.New("runtime error")
)
