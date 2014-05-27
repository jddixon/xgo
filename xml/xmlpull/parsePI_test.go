package xmlpull

// xgo/xml/xmlpull/parse_pi_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
	"io"
	"strings"
)

var _ = fmt.Print

const (
	//                ....x....1....x....2....x...."
	SIMPLE_PI         = "<?foo bar bar bardy bar ?>"
	MID_QMARK_PI      = "<?fo bar ?bar ?>"
	DASH_QMARK_AT_END = "<?f bar bar -?>"
	ENDLESS_PI        = "<?foo bar bar bardy bar"
)

func (s *XLSuite) TestSimplePI(c *C) {

	var rd1 io.Reader = strings.NewReader(SIMPLE_PI)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first  characters
	found, err := p.AcceptStr("<?")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	isPI, err := p.parsePI()
	c.Assert(err, IsNil)
	c.Assert(isPI, Equals, true)
	c.Assert(SameRunes(p.piTarget, []rune("foo")), Equals, true)
	c.Assert(SameRunes(p.piChars, []rune("bar bar bardy bar ")), Equals, true)
}

func (s *XLSuite) TestMidQmarkPI(c *C) {

	var rd1 io.Reader = strings.NewReader(MID_QMARK_PI)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 2 characters
	found, err := p.AcceptStr("<?")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	isPI, err := p.parsePI()
	c.Assert(err, IsNil)
	c.Assert(isPI, Equals, true)
	c.Assert(SameRunes(p.piTarget, []rune("fo")), Equals, true)
	c.Assert(SameRunes(p.piChars, []rune("bar ?bar ")), Equals, true)
}

func (s *XLSuite) TestDashQMarkAtEnd(c *C) {

	var rd1 io.Reader = strings.NewReader(DASH_QMARK_AT_END)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 2 characters
	found, err := p.AcceptStr("<?")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	isPI, err := p.parsePI()
	c.Assert(err, IsNil)
	c.Assert(isPI, Equals, true)
	c.Assert(SameRunes(p.piTarget, []rune("f")), Equals, true)
}

func (s *XLSuite) TestEndlessPI(c *C) {

	var rd1 io.Reader = strings.NewReader(ENDLESS_PI)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 2 characters
	found, err := p.AcceptStr("<?")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	_, err = p.parsePI()
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals,
		"processing instruction started line 1 column 3 not closed")
}
