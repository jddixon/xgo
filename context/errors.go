package context

import (
	e "errors"
)

var (
	EmptyName = e.New("name parameter is empty")
	NilValue  = e.New("value parameter is nil")
)
