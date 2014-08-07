package xmlpull

// xgo/xml/xmlpull/parseMisc_test.go

import (
	"fmt"
	"io"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
	"strings"
)

var _ = fmt.Print

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
func (s *XLSuite) setUpParser(c *C, input string) (p *Parser) {

	// DEBUG
	fmt.Printf("SET_UP_PARSER: INPUT = '%s'\n", input)
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

func (s *XLSuite) doTestParseMisc(c *C, input string,
	misc1 []*MiscItem) (p *Parser, event PullEvent) {

	p = s.setUpParser(c, input)

	//////////////////////////////////////////////////////////
	// THIS IS OUR OWN LOCAL COPY OF THE EVENT, NOT p.curEvent
	//////////////////////////////////////////////////////////
	event, err := p.NextToken()
	if err != io.EOF {
		c.Assert(err, IsNil)
	}

	lenMisc := len(misc1)
	for i := 0; i < lenMisc; i++ {
		// DEBUG
		fmt.Printf("Misc[%d/%d]: event is %s\n",
			i, lenMisc, PULL_EVENT_NAMES[event])
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
			// DEBUG
			fmt.Printf("p.text is    '%s'\n", 
				s.dumpWhiteSpace(string(p.text)))
			fmt.Printf("misc.body is '%s'\n", 
				s.dumpWhiteSpace(string(misc.body)))
			// END
			c.Assert(string(p.text), Equals, misc.body)
		}
		event, err = p.NextToken()
	}
	return
}
func (s *XLSuite) dumpWhiteSpace(sIn string) (sOut string) {
	var ss []string
	// DEBUG
	i := 0
	_ = i
	// END
	for i := 0; i < len(sIn); i++ {
		ch := sIn[i]
		ss = append(ss, fmt.Sprintf("%02x", ch))
	}
	sOut = strings.Join(ss, " ")
	return
}
	
func (s *XLSuite) TestParseMisc(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_PARSE_MISC")
	}
	rng := xr.MakeSimpleRNG()
	var misc []*MiscItem
	for {
		misc = s.createMiscItems(rng) // a small, possibly empty, slice
		if len(misc) > 0 {
			// DEBUG
			fmt.Printf("created %d MiscItem\n", len(misc))
			// END
			break
		}
	}
	p, event := s.doTestParseMisc(c, s.textFromMISlice(misc), misc)
	_, _ = p, event

	// making sure that we all agree on what white space is
	c.Assert( p.IsS('\t'), Equals, true)
	c.Assert( p.IsS('\r'), Equals, true)
	c.Assert( p.IsS('\n'), Equals, true)
	c.Assert( p.IsS(' '), Equals, true)
}
