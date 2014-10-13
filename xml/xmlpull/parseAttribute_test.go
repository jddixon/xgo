package xmlpull

// xgo/xml/xmlpull/parseAttribute_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
	//"io"
	"strings"
	u "unicode"
)

var _ = fmt.Print

// STRUCT FOR USE IN TESTING ========================================

type AttrValPair struct {
	Attr string
	Val  string
}

// Create a single randomly chosen AttrValPair.   Currently both Attr
// and Val are simple strings.  In the serialization Val will be DQUOTEd.
// For testing we require that Attr be at least one character long
// whereas Val may be zero or more characters long.

func (s *XLSuite) createAttrValPair(rng *xr.PRNG) *AttrValPair {
	var attr, val string

	attr = rng.NextFileName(8) // so 1 to 7 characters long
	for u.IsDigit(rune(attr[0])) {
		attr = rng.NextFileName(8)
	}
	val = rng.NextFileName(8) // XXX len of zero should be OK
	return &AttrValPair{Attr: attr, Val: val}
}

func (av *AttrValPair) String() (text string) {

	text = fmt.Sprintf(" %s=\"%s\"", av.Attr, av.Val)
	return
}

// UTILITY FUNCTIONS ================================================

// Returns a slice of zero or more Attributes.  Attribute names must be
// unique within the slice.
func (s *XLSuite) createAttrValPairs(rng *xr.PRNG) (pairs []*AttrValPair) {
	count := rng.Intn(4) // so 0 to 3 inclusive
	var byName = make(map[string]*AttrValPair)
	for i := 0; i < count; i++ {
		var pair *AttrValPair
		for {
			pair = s.createAttrValPair(rng)
			// attr names must be unique; values need not be
			name := pair.Attr
			if _, ok := byName[name]; ok {
				continue
			} else {
				// it's not in the map, so add it
				byName[name] = pair
				break
			}
		}
		pairs = append(pairs, pair)
	}
	return
}
func (s *XLSuite) textFromAttrValPair(pairs []*AttrValPair) string {
	var ss []string
	for i := 0; i < len(pairs); i++ {
		ss = append(ss, pairs[i].String())
	}
	// DEBUG
	fmt.Printf("  TEXT_FROM_AV_PAIRS: '%s'\n", strings.Join(ss, ""))
	// END
	return strings.Join(ss, "")
}

// UNIT TEST ========================================================

func (s *XLSuite) setupAttrParser(c *C, input string) (p *Parser) {

	// DEBUG
	fmt.Printf("SET_UP_ATTR_PARSER: INPUT = '%s'\n", input)
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

func (s *XLSuite) doTestParseAttr(c *C, input string,
	pairs []*AttrValPair) (p *Parser, event PullEvent) {

	p = s.setupAttrParser(c, input)

	//////////////////////////////////////////////////////////
	// THIS IS OUR OWN LOCAL COPY OF THE EVENT, NOT p.curEvent
	//////////////////////////////////////////////////////////
	event, err := p.NextToken()
	c.Assert(err, IsNil)
	fmt.Printf("first event:  %s\n", PULL_EVENT_NAMES[event])

	event, err = p.NextToken()
	c.Assert(err, IsNil)
	fmt.Printf("second event: %s\n", PULL_EVENT_NAMES[event])

	// We are parsing an element and the AV
	// pairs, if any, are right after the tag name

	lenAVPairs := len(pairs)
	c.Assert(p.attributeCount, Equals, uint(lenAVPairs))

	for i := 0; i < lenAVPairs; i++ {
		pair := pairs[i]
		expectedAttr := pair.Attr
		expectedVal := pair.Val
		c.Assert(p.attributeName[i], Equals, expectedAttr)
		c.Assert(p.attributeValue[i], Equals, expectedVal)

	}
	return
}

// AAA
func (s *XLSuite) TestAAAParseAVPair(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_AV_PAIR")
	}
	rng := xr.MakeSimpleRNG()
	var pairs []*AttrValPair
	for {
		pairs = s.createAttrValPairs(rng) // a small, possibly empty, slice
		if len(pairs) > 0 {
			// DEBUG
			fmt.Printf("created %d AttrValPair\n", len(pairs))
			// END
			break
		}
	}
	input := XML_DECL + DOCTYPE_DECL + "<root" + s.textFromAttrValPair(pairs) + "/>"
	p, event := s.doTestParseAttr(c, input, pairs)
	_, _ = p, event

}
