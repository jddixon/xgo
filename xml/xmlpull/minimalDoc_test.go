package xmlpull

// xgo/xml/xmlpull/minimalDocl_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
	"strings"
)

var _ = fmt.Print

const (
	XML_DECL     = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	DOCTYPE_DECL = "<!DOCTYPE document PUBLIC \"-//APACHE//DTD Documentation V2.0//EN\" \"http://forrest.apache.org/dtd/document-v20.dtd\">\n"
	PROLOG_MISC  = "<!-- this is a comment in the prolog -->\n"
	EPILOG_MISC  = "<!-- this is a comment in the epilog -->\n"
	EMPTY_ELM    = "<document/>"
)

func (s *XLSuite) doInitialParse(c *C, input string) (p *Parser) {

	// DEBUG
	fmt.Printf("PARSING: '%s'\n", input)
	// END
	rd := strings.NewReader(input)
	p, err := NewNewParser(rd)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)
	c.Assert(p.state, Equals, PRE_START_DOC)

	event, err := p.NextToken()
	c.Assert(err, IsNil)
	c.Assert(event, Equals, START_DOCUMENT)
	c.Assert(p.state, Equals, START_STATE)
	return
}

func (s *XLSuite) doParseXmlDecl(c *C, input string) (
	p *Parser, event PullEvent) {

	p = s.doInitialParse(c, input)

	event, err := p.NextToken()
	c.Assert(err, IsNil)

	c.Assert(p.xmlVersion, Equals, "1.0")
	c.Assert(p.xmlEncoding, Equals, "UTF-8")

	event, err = p.NextToken()
	c.Assert(err, IsNil)
	return
}

func (s *XLSuite) doParseBothDecl(c *C, input string) (
	p *Parser, event PullEvent) {

	p, event = s.doParseXmlDecl(c, input)
	c.Assert(event, Equals, DOCDECL)

	event, err := p.NextToken()
	c.Assert(err, IsNil)
	return
}

func (s *XLSuite) TestParseEmptyElm(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_EMPTY_ELM")
	}
	s.doInitialParse(c, EMPTY_ELM)
}

func (s *XLSuite) TestParseXmlDeclPlusElm(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_XML_DECL_PLUS_ELM")
	}
	s.doParseXmlDecl(c, XML_DECL+EMPTY_ELM)
}

func (s *XLSuite) TestParseBothDeclPlusElm(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_BOTH_DECL_PLUS_ELM")
	}
	s.doParseBothDecl(c, XML_DECL+DOCTYPE_DECL+EMPTY_ELM)
}

//func (s *XLSuite) TestParseMinimalDoc(c *C) {
//	s.doParseMinimalDoc(c, XML_DECL+DOCTYPE_DECL+PROLOG_MISC+EMPTY_ELM)
//	s.doParseMinimalDoc(c,
//		XML_DECL+DOCTYPE_DECL+PROLOG_MISC+EMPTY_ELM+EPILOG_MISC)
//}

func (s *XLSuite) TestParserConst(c *C) {
	c.Assert(PARSER_STATE_NAMES[PRE_START_DOC], Equals, "PRE_START_DOC")
	c.Assert(PARSER_STATE_NAMES[START_ROOT_SEEN], Equals, "START_ROOT_SEEN")
	c.Assert(PARSER_STATE_NAMES[PAST_END_DOC], Equals, "PAST_END_DOC")
}
