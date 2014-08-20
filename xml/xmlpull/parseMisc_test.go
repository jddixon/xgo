package xmlpull

// xgo/xml/xmlpull/parseMisc_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
	"io"
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
	body  []rune
}

var eventForMiscType = make(map[MiscType]PullEvent)

func init() {
	eventForMiscType[MISC_COMMENT] = COMMENT
	eventForMiscType[MISC_PI] = PROCESSING_INSTRUCTION
	eventForMiscType[MISC_S] = IGNORABLE_WHITESPACE
}

// Create a single randomly chosen MiscItem. If sOK it may be an S.  In any
// case it may be either a Comment or a PI.
func (s *XLSuite) createMiscItem(sOK bool, rng *xr.PRNG) *MiscItem {
	var body []rune
	var t MiscType
	if sOK {
		t = MiscType(rng.Intn(int(MISC_S) + 1))
	} else {
		t = MiscType(rng.Intn(int(MISC_S)))
	}
	switch t {
	case MISC_COMMENT:
		// The comment must not end with a dash
		for {
			body = []rune(rng.NextFileName(16)) // a quasi-random string, len < 16
			text := string(body)
			if !strings.HasSuffix(text, "-") {
				break
			}
		}
	case MISC_PI:
		body = []rune(rng.NextFileName(16)) // a quasi-random string, len < 16
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
		body = runes
	}
	return &MiscItem{_type: t, body: body}
}

func (s *XLSuite) IsS(ch rune) bool {
	return (ch == ' ') || (ch == '\t') || (ch == '\n') || (ch == '\r')
}

// Returns a slice of zero or more MiscItems.  The slice must not contain
// any S-S sequences (which are indistinguishable from a single S.
func (s *XLSuite) createMiscItems(rng *xr.PRNG) (items []*MiscItem) {
	count := rng.Intn(4) // so 0 to 3 inclusive
	lastWasS := false
	for i := 0; i < count; i++ {
		item := s.createMiscItem(true, rng) // true = S ok
		lastWasS = s.IsS(item.body[0])
		for item._type == MISC_S && lastWasS {
			item = s.createMiscItem(!lastWasS, rng)
			lastWasS = s.IsS(item.body[0])
		}
		lastWasS = item._type == MISC_S
		items = append(items, item)
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

	bodyBit := string(mi.body)
	switch mi._type {
	case MISC_COMMENT:
		text = "<!--" + bodyBit + "-->"
	case MISC_PI:
		text = "<?lang " + bodyBit + "?>"
	case MISC_S:
		text = bodyBit
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
		item := misc1[i]
		// DEBUG
		fmt.Printf("Misc[%d/%d]: event is %-20s, body is %s\n",
			i, lenMisc, PULL_EVENT_NAMES[event],
			s.dumpStrAsHex(string(item.body)))
		// END
		t := item._type
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
			c.Assert(string(p.commentChars), Equals, string(item.body))
		case MISC_PI:
			c.Assert(string(p.piChars), Equals, string(item.body))
		case MISC_S:
			// DEBUG
			fmt.Printf("p.text is    '%s'\n",
				s.dumpStrAsHex(string(p.text)))
			fmt.Printf("item.body is '%s'\n",
				s.dumpStrAsHex(string(item.body)))
			// END
			c.Assert(string(p.text), Equals, string(item.body))
		}
		event, err = p.NextToken()
	}
	return
}
func (s *XLSuite) dumpStrAsHex(sIn string) (sOut string) {
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
	c.Assert(p.IsS('\t'), Equals, true)
	c.Assert(p.IsS('\r'), Equals, true)
	c.Assert(p.IsS('\n'), Equals, true)
	c.Assert(p.IsS(' '), Equals, true)
}
