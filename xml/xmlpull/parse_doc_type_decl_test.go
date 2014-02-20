package xmlpull

// xgo/xml/xmlpull/parse_doc_decl_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

const (
	WHATEVER = "abc123"

	GENERAL_SYNTAX = "<!DOCTYPE root-element PUBLIC \"FPI\" [\"URI\"] [ \n" +
		"<!-- internal subset declarations -->\n" +
		"]>"

	COMMON_EXAMPLE = "<!DOCTYPE html PUBLIC\n" +
		"\"-//W3C//DTD XHTML 1.0 Transitional//EN\"\n" +
		"\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"

	HTML5_EXAMPLE = "<!DOCTYPE html>"
)

func (s *XLSuite) doTestParseDocTypeDecl(c *C, sample string) {

	input := sample + WHATEVER
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

	s.doTestParseDocTypeDecl(c, GENERAL_SYNTAX)
	s.doTestParseDocTypeDecl(c, COMMON_EXAMPLE)
	s.doTestParseDocTypeDecl(c, HTML5_EXAMPLE)

}