package util

import (
	e "errors"
)

var (
	NonPositiveBufferSize = e.New("non-positive buffer size")
)
