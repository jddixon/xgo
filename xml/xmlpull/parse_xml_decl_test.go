package xmlpull

// xgo/xml/xmlpull/parse_xml_decl_test.go

import (
	"fmt"
	gl "github.com/jddixon/xgo/lex"
	"io"
	. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

const (
	BASIC_DECL         = "<?xml version='1.0' ?>"
	SPACEY_BASIC_DECL  = "<?xml  version  =  '1.0'    ?>"
	DECL_WITH_ENCODING = "<?xml version='1.0' encoding='utf-8' ?>"
	STANDALONE_DECL    = "<?xml version='1.0' standalone = 'yes' ?>"
	FULL_DECL          = "<?xml version='1.0' encoding = 'utf-8' standalone = 'yes' ?>"
)

func (s *XLSuite) TestBasicDecl(c *C) {

	var rd1 io.Reader = strings.NewReader(BASIC_DECL)
	lx, err := gl.NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)
	p := &Parser{
		LexInput: *lx,
	}
	// consume the first 5 characters
	found, err := lx.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	version, encoding, standalone, err := p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(version, Equals, "1.0")
	c.Assert(standalone, Equals, false)
	c.Assert(encoding, Equals, "")
}
func (s *XLSuite) TestSpaceyBasicDecl(c *C) {

	var rd1 io.Reader = strings.NewReader(SPACEY_BASIC_DECL)
	lx, err := gl.NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)
	p := &Parser{
		LexInput: *lx,
	}
	// consume the first 5 characters
	found, err := lx.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	version, encoding, standalone, err := p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(version, Equals, "1.0")
	c.Assert(standalone, Equals, false)
	c.Assert(encoding, Equals, "")
}
func (s *XLSuite) TestDeclWithEncoding(c *C) {

	var rd1 io.Reader = strings.NewReader(DECL_WITH_ENCODING)
	lx, err := gl.NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)
	p := &Parser{
		LexInput: *lx,
	}
	// consume the first 5 characters
	found, err := lx.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	version, encoding, standalone, err := p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(version, Equals, "1.0")
	c.Assert(standalone, Equals, false)
	c.Assert(encoding, Equals, "utf-8")
}
func (s *XLSuite) TestStandaloneDecl(c *C) {

	var rd1 io.Reader = strings.NewReader(STANDALONE_DECL)
	lx, err := gl.NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)
	p := &Parser{
		LexInput: *lx,
	}
	// consume the first 5 characters
	found, err := lx.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	version, encoding, standalone, err := p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(version, Equals, "1.0")
	c.Assert(standalone, Equals, true)
	c.Assert(encoding, Equals, "")
}
func (s *XLSuite) TestFullDecl(c *C) {

	var rd1 io.Reader = strings.NewReader(FULL_DECL)
	lx, err := gl.NewLexInput(rd1, "") // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(lx, NotNil)
	p := &Parser{
		LexInput: *lx,
	}
	// consume the first 5 characters
	found, err := lx.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	version, encoding, standalone, err := p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(version, Equals, "1.0")
	c.Assert(standalone, Equals, true)
	c.Assert(encoding, Equals, "utf-8")
}
