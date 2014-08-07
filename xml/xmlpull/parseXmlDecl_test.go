package xmlpull

// xgo/xml/xmlpull/parse_xml_decl_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
	"io"
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

	if VERBOSITY > 0 {
		fmt.Println("TEST_BASIC_DECL")
	}
	var rd1 io.Reader = strings.NewReader(BASIC_DECL)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 5 characters
	found, err := p.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(p.xmlVersion, Equals, "1.0")
	c.Assert(p.xmlDeclStandalone, Equals, false)
	c.Assert(p.xmlEncoding, Equals, "")
}
func (s *XLSuite) TestSpaceyBasicDecl(c *C) {

	if VERBOSITY > 0 {
		fmt.Println("TEST_SPACEY_BASIC_DECL")
	}
	var rd1 io.Reader = strings.NewReader(SPACEY_BASIC_DECL)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 5 characters
	found, err := p.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(p.xmlVersion, Equals, "1.0")
	c.Assert(p.xmlDeclStandalone, Equals, false)
	c.Assert(p.xmlEncoding, Equals, "")
}
func (s *XLSuite) TestDeclWithEncoding(c *C) {

	if VERBOSITY > 0 {
		fmt.Println("TEST_DECL_WITH_ENCODING")
	}
	var rd1 io.Reader = strings.NewReader(DECL_WITH_ENCODING)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 5 characters
	found, err := p.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(p.xmlVersion, Equals, "1.0")
	c.Assert(p.xmlDeclStandalone, Equals, false)
	c.Assert(p.xmlEncoding, Equals, "utf-8")
}
func (s *XLSuite) TestStandaloneDecl(c *C) {

	if VERBOSITY > 0 {
		fmt.Println("TEST_STANDALONE_DECL")
	}
	var rd1 io.Reader = strings.NewReader(STANDALONE_DECL)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 5 characters
	found, err := p.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(p.xmlVersion, Equals, "1.0")
	c.Assert(p.xmlDeclStandalone, Equals, true)
	c.Assert(p.xmlEncoding, Equals, "")
}
func (s *XLSuite) TestFullDecl(c *C) {

	if VERBOSITY > 0 {
		fmt.Println("TEST_FULL_DECL")
	}
	var rd1 io.Reader = strings.NewReader(FULL_DECL)
	p, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first 5 characters
	found, err := p.AcceptStr("<?xml")
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	err = p.parseXmlDecl()
	c.Assert(err, IsNil)
	c.Assert(p.xmlVersion, Equals, "1.0")
	c.Assert(p.xmlDeclStandalone, Equals, true)
	c.Assert(p.xmlEncoding, Equals, "utf-8")
}
