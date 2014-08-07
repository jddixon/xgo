package xmlpull

// xgo/xml/xmlpull/parseDocTypeDecl_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
	"strings"
)

var _ = fmt.Print

const (
	WHATEVER   = "abc123"
	OPEN_DECL  = "<!DOCTYPE "
	CLOSE_DECL = ">"

	GENERAL_SYNTAX = OPEN_DECL + "root-element PUBLIC \"FPI\" [\"URI\"] [ \n" +
		"<!-- internal subset declarations -->\n]"

	COMMON_EXAMPLE = "html PUBLIC\n" +
		"\"-//W3C//DTD XHTML 1.0 Transitional//EN\"\n" +
		"\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\""

	HTML5_EXAMPLE = "html"
)

func (s *XLSuite) doTestParseDocTypeDecl(c *C, sample string) {

	input := OPEN_DECL + sample + CLOSE_DECL + WHATEVER
	rd := strings.NewReader(input)
	p, err := NewNewParser(rd)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	// consume the first three characters
	err = p.ExpectStr("<!D")
	c.Assert(err, IsNil)

	err = p.parseDocTypeDecl()
	c.Assert(err, IsNil)
	c.Assert(p.docTypeDecl, Equals, sample)

	err = p.ExpectStr(WHATEVER)
	c.Assert(err, IsNil)
}
func (s *XLSuite) TestParseDocTypeDecl(c *C) {

	if VERBOSITY > 0 {
		fmt.Println("TEST_PARSE_DOC_TYPE_DECL")
	}
	s.doTestParseDocTypeDecl(c, GENERAL_SYNTAX)
	s.doTestParseDocTypeDecl(c, COMMON_EXAMPLE)
	s.doTestParseDocTypeDecl(c, HTML5_EXAMPLE)

}
