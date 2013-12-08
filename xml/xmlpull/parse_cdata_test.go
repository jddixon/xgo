package xmlpull

// xgo/xml/xmlpull/parse_cdata_test.go

import (
	"fmt"
	"io"
	. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

const (
	//                  ....x....1....x....2....x."
	SIMPLE_CDSECT      = "<![CDATA[ foo foo foody foo ]]>"
	TWO_BRACKET_CDSECT = "<![CDATA[ foo foo ]] foody foo ]]>"
	ENDLESS_CDSECT     = "<![CDATA[ foo foo foody foo"
)

func (s *XLSuite) TestSimpleCDSect(c *C) {

	var rd1 io.Reader = strings.NewReader(SIMPLE_CDSECT)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 3 chars (rest of prefix will be discarded in parse)
	found, err := p.AcceptStr("<![")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseCDSect()
	c.Assert(err, IsNil)
	c.Assert(p.cDataChars, Equals, " foo foo foody foo ")
}

func (s *XLSuite) TestTwoBracketCDSect(c *C) {

	var rd1 io.Reader = strings.NewReader(TWO_BRACKET_CDSECT)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 3 chars (rest of prefix will be discarded in parse)
	found, err := p.AcceptStr("<![")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseCDSect()
	c.Assert(err, IsNil)
	c.Assert(p.cDataChars, Equals, " foo foo ]] foody foo ")
}

func (s *XLSuite) TestEndlessCDSect(c *C) {

	var rd1 io.Reader = strings.NewReader(ENDLESS_CDSECT)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 3 chars (rest of prefix will be discarded in parse)
	found, err := p.AcceptStr("<![")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseCDSect()
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "cData started line 1 column 9 not closed")
}
