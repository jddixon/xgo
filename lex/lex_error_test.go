package lex

// xgo/lex/lex_input_test.go

import (
	"fmt"
	"io"
	. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

func (s *XLSuite) TestLexError(c *C) {

	const (
		ERR_MSG_1 = "test msg 1"
		ERR_MSG_2 = "second test msg"
		ERR_MSG_3 = "error msg 3"
	)
	var rd1 io.Reader = strings.NewReader("abc\ndef\nversion 97.1 ")
	lx, err := NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)

	err = lx.ExpectStr("abc")
	c.Assert(err, IsNil)

	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 3)

	e1 := lx.NewError(ERR_MSG_1)
	expected := fmt.Sprintf("line %d col %d: %s",
		lx.LineNo(), lx.ColNo(), ERR_MSG_1)
	c.Assert(e1.Error(), Equals, expected)

	nl, err := lx.NextCh()
	c.Assert(err, IsNil)
	c.Assert(nl, Equals, '\n')
	c.Assert(lx.LineNo(), Equals, 2)
	c.Assert(lx.ColNo(), Equals, 0)

	e2 := lx.NewError(ERR_MSG_2)
	expected = fmt.Sprintf("line %d col %d: %s",
		lx.LineNo(), lx.ColNo(), ERR_MSG_2)
	c.Assert(e2.Error(), Equals, expected)

	err = lx.ExpectStr("def")
	c.Assert(err, IsNil)
	c.Assert(lx.LineNo(), Equals, 2)
	c.Assert(lx.ColNo(), Equals, 3)

	err = lx.ExpectCh('\n')
	c.Assert(err, IsNil)
	c.Assert(lx.LineNo(), Equals, 3)
	c.Assert(lx.ColNo(), Equals, 0)

	err = lx.ExpectStr("version 97.1 ")
	c.Assert(err, IsNil)
	c.Assert(lx.LineNo(), Equals, 3)
	c.Assert(lx.ColNo(), Equals, 13)

	e3 := lx.NewError(ERR_MSG_3)
	expected = fmt.Sprintf("line %d col %d: %s",
		lx.LineNo(), lx.ColNo(), ERR_MSG_3)
	c.Assert(e3.Error(), Equals, expected)

}
