package lex

// xgo/lex/ji_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
	u "unicode"
	"unicode/utf8"
)

var _ = fmt.Print

var (
	NIHONGO      = "This is \u65e5\u672c\u8a93"
	NIHONGO_RUNE = 11
	NIHONGO_BYTE = 17
	TOKYO        = "私たちは、約3年半、東京の若松町に住んでいました。"
	TOKYO_RUNE   = 25
	TOKYO_BYTE   = 73
	HIRAGANA     []rune
	KATAKANA     []rune
)

// We want this test to run first, to populate HIRAGANA and KATAKANA.
//
func (s *XLSuite) TestAAA(c *C) {

	// 3 bytes per kanji
	c.Assert(len(NIHONGO), Equals, NIHONGO_BYTE)
	c.Assert(utf8.RuneCountInString(NIHONGO), Equals, NIHONGO_RUNE)
	fmt.Printf("THE ENTIRE STRING:     '%s'\n", NIHONGO)
	fmt.Printf("A CHARACTER AT A TIME: '")

	// This prints out each Unicode character, each rune, as expected.
	for _, ch := range NIHONGO {
		fmt.Printf("%c", ch)
	}
	fmt.Println()

	c.Assert(utf8.RuneCountInString(TOKYO), Equals, TOKYO_RUNE)
	// 3 bytes per kanji but 1 for the digit '3'
	c.Assert(len(TOKYO), Equals, TOKYO_BYTE)

	for n, ch := range TOKYO {
		// no spaces in Japanese ;-)
		c.Assert(u.IsSpace(ch), Equals, false)
		fmt.Printf("%#U  ", ch)
		if ch == 0x0033 {
			fmt.Print(" ") // %#U puts a space after ji
		} else if n == 15 || n == 31 || n == 49 || n == 67 {
			fmt.Println()
		}
	}
	fmt.Println()

	// hiragana 3041 through 3093
	fmt.Println("hiragana")
	for r := rune(0x3041); r < rune(0x3094); r++ {
		fmt.Printf("%c", r)
		if r == 0x304a || r == 0x3054 || r == 0x305e || r == 0x3069 ||
			r == 0x306e || r == 0x307d || r == 0x3082 || r == 0x3088 ||
			r == 0x308d {

			fmt.Println()
		}
		HIRAGANA = append(HIRAGANA, r)
	}
	fmt.Println()
	strH := string(HIRAGANA)
	// hiragana are 3 bytes each
	c.Assert(3*len(HIRAGANA), Equals, len(strH))

	// hatakana 30a1 through 30ef
	fmt.Println("katakana")
	for r := rune(0x30a1); r < 0x30f0; r++ {
		fmt.Printf("%c", r)
		if r == 0x30aa || r == 0x30b4 || r == 0x30be || r == 0x30c9 ||
			r == 0x30ce || r == 0x30d4 || r == 0x30dc || r == 0x30e2 ||
			r == 0x30e8 || r == 0x30ed {

			fmt.Println()
		}
		KATAKANA = append(KATAKANA, r)
	}
	fmt.Println()
	strK := string(KATAKANA)
	c.Assert(3*len(KATAKANA), Equals, len(strK))
}
