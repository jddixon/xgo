package xmlpull

// xgo/xml/xmlpull/minimalDoc_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
	"io"
	"strings"
)

var _ = fmt.Print

const (
	XML_DECL     = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	DOCTYPE_DECL = "<!DOCTYPE document PUBLIC \"-//APACHE//DTD Documentation V2.0//EN\" \"http://forrest.apache.org/dtd/document-v20.dtd\">"
	PROLOG_MISC  = "<!-- this is a comment in the prolog -->\n"
	EPILOG_MISC  = "<!-- this is a comment in the epilog -->\n"
	EMPTY_ELM    = "<root/>"
)

func (s *XLSuite) doInitialParse(c *C, input string) (p *Parser) {

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

// Parse a sequence: XmlDecl Misc* (doctypedecl Misc*) EmptyElement Misc*
func (s *XLSuite) doParseBothDecl(c *C, input string) (
	p *Parser, event PullEvent) {

	var err error
	p, event = s.doParseXmlDecl(c, input)

	// we have seen the XmlDecl; now allow zero or more Misc but
	// requuire doctypedecl
	for event != DOCDECL {
		if event == IGNORABLE_WHITESPACE || event == PROCESSING_INSTRUCTION ||
			event == COMMENT {

			event, err = p.NextToken()
		} else {
			fmt.Printf("expected DOCDECL but got %s\n", PULL_EVENT_NAMES[event])
		}
	}
	c.Assert(event, Equals, DOCDECL)

	event, err = p.NextToken()
	c.Assert(err, IsNil)

	// Allow zero or more Misc but require the EmptyElement
	for event != START_TAG {
		if event == IGNORABLE_WHITESPACE || event == PROCESSING_INSTRUCTION ||
			event == COMMENT {

			event, err = p.NextToken()
		} else {
			fmt.Printf("expected START_TAG but got %s\n", PULL_EVENT_NAMES[event])
			break // XXX SHOULD FAIL
		}
	}
	c.Assert(err, IsNil)
	c.Assert(event, Equals, START_TAG)

	// allow any arbitrary number of Misc
	event, err = p.NextToken()
	// DEBUG
	fmt.Printf("DoParseBothDecl: NextToken => event %s, err %v\n",
		PULL_EVENT_NAMES[event], err)
	// END
	c.Assert(err == nil || err == io.EOF, Equals, true)
	fmt.Printf("err: %v\n", err)
	for event == IGNORABLE_WHITESPACE ||
		event == PROCESSING_INSTRUCTION || event == COMMENT {

		event, err = p.NextToken()
		// DEBUG
		fmt.Printf("DoParseBothDecl: NextToken => event %s, err %v\n",
			PULL_EVENT_NAMES[event], err)
		// END
		c.Assert(err == nil || err == io.EOF, Equals, true)
		fmt.Printf("err: %v\n", err)
	}
	c.Assert(event, Equals, END_DOCUMENT)
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

// Parse an XmlDecl followed by an (empty) element followed by Misc
func (s *XLSuite) TestParseXmlDeclPlusElmPlusMisc(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_XML_DECL_PLUS_ELM_PLUS_MISC")
	}
	rng := xr.MakeSimpleRNG()
	misc1 := s.createMiscItems(rng) // a small, possibly empty, slice
	miscN := s.createMiscItems(rng) // a small, possibly empty, slice
	s.doParseXmlDeclWithMisc(c, XML_DECL+s.textFromMISlice(misc1)+
		EMPTY_ELM+s.textFromMISlice(miscN), misc1)
}

func (s *XLSuite) TestParseBothDeclPlusElm(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_BOTH_DECL_PLUS_ELM")
	}
	s.doParseBothDecl(c, XML_DECL+DOCTYPE_DECL+EMPTY_ELM)
}

// Parse an XmlDecl followed a DocDecl followed by an (empty) element followed
// by Misc
// [1]  document ::= prolog element Misc*
// [22] prolog ::= XMLDecl? Misc* (doctypedecl Misc*)?
//
func (s *XLSuite) TestParseBothDeclPlusElmPlusMisc(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_BOTH_DECL_PLUS_ELM_PLUS_MISC")
	}
	rng := xr.MakeSimpleRNG()
	misc1 := s.createMiscItems(rng) // a small, possibly empty, slice
	misc2 := s.createMiscItems(rng) // a small, possibly empty, slice
	miscN := s.createMiscItems(rng) // a small, possibly empty, slice
	s.doParseBothDecl(c,
		XML_DECL+s.textFromMISlice(misc1)+
			DOCTYPE_DECL+s.textFromMISlice(misc2)+
			EMPTY_ELM+s.textFromMISlice(miscN))
}

// Simple test that constants and their string representations agree.
//
func (s *XLSuite) TestParserConst(c *C) {
	c.Assert(PARSER_STATE_NAMES[PRE_START_DOC], Equals, "PRE_START_DOC")
	c.Assert(PARSER_STATE_NAMES[START_ROOT_SEEN], Equals, "START_ROOT_SEEN")
	c.Assert(PARSER_STATE_NAMES[PAST_END_DOC], Equals, "PAST_END_DOC")
}
