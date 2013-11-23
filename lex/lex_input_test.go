package lex

// xgo/lex/lex_input_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
)
var _ = fmt.Print

func (s *XLSuite) TestS (c *C) {

	whitespace := []byte(" \n\r\t")
	for i := 0; i < len(whitespace) ; i++ {
		c.Assert(IsS(whitespace[i]), Equals, true)
	}
}
