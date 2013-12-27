package md

// xgo/md/blocks_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	"io"
	. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

// -- HEADER --------------------------------------------------------

func (s *XLSuite) TestHeader(c *C) {
	rng := xr.MakeSimpleRNG()

	titles := make([]string, 6)
	for i := 0; i < 6; i++ {
		titles[i] = rng.NextFileName(16)
	}
	hashes := make([]string, 6)
	for j := 0; j < 6; j++ {
		if j == 0 {
			hashes[j] = "#"
		} else {
			hashes[j] = hashes[j-1] + "#"
		}
	}

	lines := make([]string, 6)
	for i := 0; i < 6; i++ {
		lines[i] = hashes[i] + titles[i]
		// half get trailing hashes
		if i%2 == 0 {
			lines[i] += hashes[i]
		}
	}
	input := strings.Join(lines, "\n")

	var rd io.Reader = strings.NewReader(input)
	p, err := NewParser(rd)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	for i := 0; i < 6; i++ {
		// test parse -------------------------------------
		line, err := p.readLine()
		if i < 5 {
			c.Assert(err, IsNil)
		} else {
			c.Assert(err, Equals, io.EOF)
		}
		c.Assert(line, NotNil)
		b, err := line.parseHeader()
		c.Assert(err, IsNil)
		c.Assert(b, NotNil)
		h := b.(*Header)
		c.Assert(h.n, Equals, i+1)
		c.Assert(string(h.runes), Equals, string(titles[i]))

		// test serialization -----------------------------
		ser := string(b.Get())
		expected := fmt.Sprintf("<h%d>%s</h%d>", i+1, titles[i], i+1)
		c.Assert(ser, Equals, expected)
	}
}

// -- HRULE ---------------------------------------------------------
// XXX Tests only strings expected to succeed.
func (s *XLSuite) doTestHRule(c *C, char rune, rng *xr.PRNG) {
	var text []rune
	runs := 3 + rng.Intn(3) // from three to five characters
	for i := 0; i < runs; i++ {
		text = append(text, char)
		if i < runs-1 {
			spaces := rng.Intn(3)
			for j := 0; j < spaces; j++ {
				text = append(text, ' ')
			}
		}
	}
	text = append(text, '\r')
	input := string(text)
	var rd io.Reader = strings.NewReader(input)
	p, err := NewParser(rd)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)

	line, err := p.readLine()
	c.Assert(err, IsNil)
	b, err := line.parseHRule()
	c.Assert(err, IsNil)
	c.Assert(b, NotNil)
	h := b.(*HRule)

	// test serialization -----------------------------
	ser := string(h.Get())
	c.Assert(ser, Equals, "<hr />")
}
func (s *XLSuite) TestHRule(c *C) {
	rng := xr.MakeSimpleRNG()

	s.doTestHRule(c, '-', rng)
	s.doTestHRule(c, '*', rng)
	s.doTestHRule(c, '_', rng)
}