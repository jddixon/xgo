package om

import (
	e "errors"
)

var (
	EmptyOtherList   = e.New("other list is empty")
	GraphCycleError  = e.New("graph cycle error")
	IndexOutOfBounds = e.New("index out of bounds")
	NilAttr          = e.New("nil Attr argument")
	NilHolder        = e.New("nil Holder argument")
	NilNamespace     = e.New("nil namespace argument")
	NilNode          = e.New("nil Node argument")
	RuntimeError     = e.New("runtime error")
	SettingDocsDoc   = e.New("setting document's document")
)
