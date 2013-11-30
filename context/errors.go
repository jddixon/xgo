package context

import (
	e "errors"
)

var (
	// The serialiation has zero length.
	EmptySerialization        = e.New("empty serialization")
	EmptyName                 = e.New("name parameter is empty")
	IllFormedCtxSerialization = e.New("ll-formed context serialization")
	NilValue                  = e.New("value parameter is nil")
)
