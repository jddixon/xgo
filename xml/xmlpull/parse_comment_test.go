package xmlpull

// xgo/xml/xmlpull/parse_comment_test.go

import (
	"fmt"
	"io"
	. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

const (
	//                     ....x....1....x....2....x."
	SIMPLE_COMMENT    = "<!-- foo foo foody foo -->"
	TWO_DASH_COMMENT  = "<!-- foo -- foo -->"
	THREE_DASH_AT_END = "<!-- foo foo --->"
	ENDLESS_COMMENT   = "<!-- foo foo foody foo"
)

func (s *XLSuite) TestSimpleComment(c *C) {

	var rd1 io.Reader = strings.NewReader(SIMPLE_COMMENT)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 3 characters (fourth will be discarded in parse)
	found, err := p.AcceptStr("<!-")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseComment()
	c.Assert(err, IsNil)
	c.Assert(p.commentChars, Equals, " foo foo foody foo ")
}

func (s *XLSuite) TestTwoDashComment(c *C) {

	var rd1 io.Reader = strings.NewReader(TWO_DASH_COMMENT)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 3 characters (fourth will be discarded in parse)
	found, err := p.AcceptStr("<!-")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseComment()
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals,
		"line 1 col 12: cannot have two dashes within comment")

}

func (s *XLSuite) TestThreeDashAtEnd(c *C) {

	var rd1 io.Reader = strings.NewReader(THREE_DASH_AT_END)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 3 characters (fourth will be discarded in parse)
	found, err := p.AcceptStr("<!-")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseComment()
	c.Assert(err, NotNil)
}

func (s *XLSuite) TestEndlessComment(c *C) {

	var rd1 io.Reader = strings.NewReader(ENDLESS_COMMENT)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 3 characters (fourth will be discarded in parse)
	found, err := p.AcceptStr("<!-")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseComment()
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "comment started line 1 column 4 not closed")
}
