package om

import (
	e "errors"
)

var (
	EmptyNamespace   = e.New("empty namespace argument")
	EmptyOtherList   = e.New("other list is empty")
	GraphCycleError  = e.New("graph cycle error")
	IndexOutOfBounds = e.New("index out of bounds")
	NilAttr          = e.New("nil Attr argument")
	NilHolder        = e.New("nil Holder argument")
	NilNode          = e.New("nil Node argument")
	NodeListNotEmpty = e.New("NodeList is not empty")
	RuntimeError     = e.New("runtime error")
	SettingDocsDoc   = e.New("setting document's document")
)
