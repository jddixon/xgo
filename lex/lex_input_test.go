package lex

// xgo/lex/lex_input_test.go

import (
	"fmt"
	"io"
	. "launchpad.net/gocheck"
	"strings"
	u "unicode"
	//"unicode/utf8"
)

var _ = fmt.Print

func (s *XLSuite) TestS(c *C) {
	whitespace := " \n\r\t"
	for _, r := range whitespace {
		c.Assert(u.IsSpace(r), Equals, true)
	}
}

func (s *XLSuite) TestStringReader(c *C) {

	var rd1 io.Reader = strings.NewReader(TOKYO)
	var runes []rune
	lx, err := NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)

	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 0)

	r, err := lx.NextCh()
	c.Assert(err, IsNil)
	c.Assert(r, Equals, rune(0x79c1))
	runes = append(runes, r)
	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 1)

	lx.SkipS() // exercises lx.pushBack()

	r, err = lx.NextCh()
	c.Assert(err, IsNil)
	c.Assert(r, Equals, rune(0x305f))
	runes = append(runes, r)
	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 2)

	_, _ = rd1, runes

}

func (s *XLSuite) TestEnglishReader(c *C) {

	//               ....x....1....x....2....x....3....x....4..
	const ENGLISH = "This    is a   test \nof many things   !  "
	var rd1 io.Reader = strings.NewReader(ENGLISH)
	var runes []rune
	lx, err := NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)

	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 0)

	r, err := lx.NextCh()
	c.Assert(err, IsNil)
	c.Assert(r, Equals, rune('T'))
	runes = append(runes, r)
	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 1)

	lx.SkipS() // exercises lx.pushBack()

	r, err = lx.NextCh()
	c.Assert(err, IsNil)
	c.Assert(r, Equals, rune('h'))
	runes = append(runes, r)
	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 2)

	r, err = lx.NextCh()
	c.Assert(err, IsNil)
	c.Assert(r, Equals, rune('i'))
	runes = append(runes, r)
	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 3)

	r, err = lx.NextCh()
	c.Assert(err, IsNil)
	c.Assert(r, Equals, rune('s'))
	runes = append(runes, r)
	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 4)

	err = lx.ExpectS() // skips 4 spaces
	c.Assert(err, IsNil)
	lx.SkipS() // redundant, of course

	r, err = lx.NextCh()
	c.Assert(err, IsNil)
	c.Assert(r, Equals, rune('i'))
	runes = append(runes, r)
	c.Assert(lx.LineNo(), Equals, 1)
	c.Assert(lx.ColNo(), Equals, 9)

	// XXX MOVE ON UP TO THE NEWLINE, PLEASE

}

func (s *XLSuite) TestExpectStr(c *C) {

	var rd1 io.Reader = strings.NewReader("version 97.1 ")
	lx, err := NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)

	err = lx.ExpectStr("version")
	c.Assert(err, IsNil)
	err = lx.ExpectS()
	c.Assert(err, IsNil)
	err = lx.ExpectStr("97.1")
	c.Assert(err, IsNil)
	err = lx.ExpectS()
	c.Assert(err, IsNil)

	rd1 = strings.NewReader("verxion 97.1 ")
	lx, err = NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)

	err = lx.ExpectStr("version")
	c.Assert(err, NotNil)

	expectedMsg := "expected 's' in 'version', found 'x'"
	c.Assert(err.Error(), Equals, expectedMsg)
}
func (s *XLSuite) TestAcceptStr(c *C) {

	var found bool
	var rd1 io.Reader = strings.NewReader("version 97.1 ")
	lx, err := NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)

	found, err = lx.AcceptStr("version")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	err = lx.ExpectS()
	c.Assert(err, IsNil)
	found, err = lx.AcceptStr("97.1")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	err = lx.ExpectS()
	c.Assert(err, IsNil)

	rd2 := strings.NewReader("verxion 97.1 ")
	lx2, err := NewLexInput(rd2, "")
	c.Assert(err, IsNil)
	c.Assert(lx2, NotNil)

	found, err = lx2.AcceptStr("version")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, false)

	found, err = lx2.AcceptStr("verxion")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	err = lx2.ExpectS()
	c.Assert(err, IsNil)
	found, err = lx2.AcceptStr("97.1")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	err = lx2.ExpectS()
	c.Assert(err, IsNil)
}
