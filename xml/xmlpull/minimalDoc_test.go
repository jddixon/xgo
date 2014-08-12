package xmlpull

// xgo/xml/xmlpull/minimalDoc_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
	"strings"
)

var _ = fmt.Print

const (
	XML_DECL     = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	DOCTYPE_DECL = "<!DOCTYPE document PUBLIC \"-//APACHE//DTD Documentation V2.0//EN\" \"http://forrest.apache.org/dtd/document-v20.dtd\">"
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

	return
}

func (s *XLSuite) doParseXmlDeclWithMisc(c *C, input string,
	misc1 []*MiscItem) (p *Parser, event PullEvent) {

	p = s.doInitialParse(c, input)

	//////////////////////////////////////////////////////////
	// THIS IS OUR OWN LOCAL COPY OF THE EVENT, NOT p.curEvent
	//////////////////////////////////////////////////////////
	event, err := p.NextToken()
	c.Assert(err, IsNil)

	c.Assert(p.xmlVersion, Equals, "1.0")
	c.Assert(p.xmlEncoding, Equals, "UTF-8")

	for i := 0; i < len(misc1); i++ {
		// DEBUG
		fmt.Printf("Misc[%d]: event is %s\n",
			i, PULL_EVENT_NAMES[event])
		// END
		misc := misc1[i]
		t := misc._type
		if event != eventForMiscType[t] {
			fmt.Printf(
				"expected event %s for misc type %s but event is %s\n",
				PULL_EVENT_NAMES[eventForMiscType[t]],
				MiscTypeNames[t],
				PULL_EVENT_NAMES[event])
		}
		c.Assert(event, Equals, eventForMiscType[t])
		switch t {
		case MISC_COMMENT:
			c.Assert(string(p.commentChars), Equals, string(misc.body))
		case MISC_PI:
			c.Assert(string(p.piChars), Equals, string(misc.body))
		case MISC_S:
			c.Assert(string(p.text), Equals, string(misc.body))
		}
		event, err = p.NextToken()
	}
	return
}

func (s *XLSuite) doParseBothDecl(c *C, input string) (
	p *Parser, event PullEvent) {

	p, event = s.doParseXmlDecl(c, input)
	// DEBUG
	if event != DOCDECL {
		fmt.Printf("expected DOCDECL but got %s\n", PULL_EVENT_NAMES[event])
	}
	// END
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

func (s *XLSuite) TestParseXmlDeclPlusElmPlusMisc(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_XML_DECL_PLUS_ELM_PLUS_MISC")
	}
	rng := xr.MakeSimpleRNG()
	misc1 := s.createMiscItems(rng) // a small, possibly empty, slice
	miscN := s.createMiscItems(rng) // a small, possibly empty, slice

	_, _ = misc1, miscN

	s.doParseXmlDeclWithMisc(c, XML_DECL+s.textFromMISlice(misc1)+
		EMPTY_ELM+s.textFromMISlice(miscN), misc1)
} // GEEP

func (s *XLSuite) TestParseBothDeclPlusElm(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_BOTH_DECL_PLUS_ELM")
	}
	s.doParseBothDecl(c, XML_DECL+DOCTYPE_DECL+EMPTY_ELM)
}

func (s *XLSuite) TestParseBothDeclPlusElmPlusMisc(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_BOTH_DECL_PLUS_ELM_PLUS_MISC")
	}
	rng := xr.MakeSimpleRNG()
	misc1 := s.createMiscItems(rng) // a small, possibly empty, slice
	misc2 := s.createMiscItems(rng) // a small, possibly empty, slice
	miscN := s.createMiscItems(rng) // a small, possibly empty, slice

	// WORKING HERE
	_, _, _ = misc1, misc2, miscN

	s.doParseBothDecl(c, XML_DECL+DOCTYPE_DECL+EMPTY_ELM)
}

// Simple test that constants and their string representations agree.
//
func (s *XLSuite) TestParserConst(c *C) {
	c.Assert(PARSER_STATE_NAMES[PRE_START_DOC], Equals, "PRE_START_DOC")
	c.Assert(PARSER_STATE_NAMES[START_ROOT_SEEN], Equals, "START_ROOT_SEEN")
	c.Assert(PARSER_STATE_NAMES[PAST_END_DOC], Equals, "PAST_END_DOC")
}
