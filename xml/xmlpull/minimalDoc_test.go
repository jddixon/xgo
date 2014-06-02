package xmlpull

// xgo/xml/xmlpull/minimalDocl_test.go

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

type MiscType int

const (
	MISC_COMMENT = MiscType(iota)
	MISC_PI
	MISC_S
)

var MiscTypeNames = []string{
	"COMMENT", "PI", "S",
}

type MiscItem struct {
	_type MiscType
	body  string
}

var eventForMiscType = make(map[MiscType]PullEvent)

func init() {
	eventForMiscType[MISC_COMMENT] = COMMENT
	eventForMiscType[MISC_PI] = PROCESSING_INSTRUCTION
	eventForMiscType[MISC_S] = IGNORABLE_WHITESPACE
}
func (s *XLSuite) createMiscItem(rng *xr.PRNG) *MiscItem {
	var body string
	t := MiscType(rng.Intn(int(MISC_S) + 1))
	switch t {
	case MISC_COMMENT:
		fallthrough
	case MISC_PI:
		body = rng.NextFileName(16) // a quasi-random string, len < 16
	case MISC_S:
		var runes []rune
		count := 1 + rng.Intn(3) // 1 to 3 inclusive
		for i := 0; i < count; i++ {
			kind := rng.Intn(4) // 0 to 3 inclusive
			switch kind {
			case 0:
				runes = append(runes, '\t')
			case 1:
				runes = append(runes, '\n')
			case 2:
				runes = append(runes, '\r')
			case 3:
				runes = append(runes, ' ')
			}
		}
		body = string(runes)
	}
	// DEBUG
	fmt.Printf("  CREATED MISC: %-7s '%s'\n", MiscTypeNames[t], body)
	// END
	return &MiscItem{_type: t, body: body}
}
func (s *XLSuite) createMiscItems(rng *xr.PRNG) (items []*MiscItem) {
	count := rng.Intn(4) // so 0 to 3 inclusive
	for i := 0; i < count; i++ {
		items = append(items, s.createMiscItem(rng))
	}
	return
}
func (s *XLSuite) textFromMISlice(items []*MiscItem) string {
	var ss []string
	for i := 0; i < len(items); i++ {
		ss = append(ss, items[i].String())
	}
	// DEBUG
	fmt.Printf("  TEXT_FROM_SLICES: '%s'\n", strings.Join(ss, ""))
	// END
	return strings.Join(ss, "")
}
func (mi *MiscItem) String() (text string) {
	switch mi._type {
	case MISC_COMMENT:
		text = "<!--" + mi.body + "-->"
	case MISC_PI:
		text = "<?lang " + mi.body + "?>"
	case MISC_S:
		text = mi.body
	}
	return
}
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
			c.Assert(string(p.commentChars), Equals, misc.body)
		case MISC_PI:
			c.Assert(string(p.piChars), Equals, misc.body)
		case MISC_S:
			c.Assert(string(p.text), Equals, misc.body)
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

// aaa effectively comments out this test
func (s *XLSuite) aaaTestParseEmptyElm(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_EMPTY_ELM")
	}
	s.doInitialParse(c, EMPTY_ELM)
}

// aaa effectively comments out this test
func (s *XLSuite) aaaTestParseXmlDeclPlusElm(c *C) {
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

	// WORKING HERE
	_, _ = misc1, miscN

	s.doParseXmlDeclWithMisc(c, XML_DECL+s.textFromMISlice(misc1)+
		EMPTY_ELM+s.textFromMISlice(miscN), misc1)
}

// aaa effectively comments out this test
func (s *XLSuite) aaaTestParseBothDeclPlusElm(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_BOTH_DECL_PLUS_ELM")
	}
	s.doParseBothDecl(c, XML_DECL+DOCTYPE_DECL+EMPTY_ELM)
}

// aaa effectively comments out this test
func (s *XLSuite) aaaTestParseBothDeclPlusElmPlusMisc(c *C) {
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
