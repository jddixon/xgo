package lex

// xgo/lex/lex_input_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
	//"strings"
	u "unicode"
	"unicode/utf8"
)

var _ = fmt.Print

func (s *XLSuite) TestS(c *C) {

	whitespace := " \n\r\t"
	for _, r := range whitespace {
		c.Assert(u.IsSpace(r), Equals, true)
	}

	nihongo := "This is \u65e5\u672c\u8a93"
	// 3 bytes per kanji
	c.Assert(len(nihongo), Equals, 17)
	c.Assert(utf8.RuneCountInString(nihongo), Equals, 11)
	fmt.Printf("THE ENTIRE STRING:     '%s'\n", nihongo)
	fmt.Printf("A CHARACTER AT A TIME: '")

	// This prints out each Unicode character, each rune, as expected.
	for _, ch := range nihongo {
		fmt.Printf("%c", ch)
	}
	fmt.Println("'")

	tokyo := "私たちは、約3年半、東京の若松町に住んでいました。"
	c.Assert(utf8.RuneCountInString(tokyo), Equals, 25)
	// 3 bytes per kanji less two
	c.Assert(len(tokyo), Equals, 73)

	for _, ch := range tokyo {
		// no spaces in Japanese ;-)
		c.Assert(u.IsSpace(ch), Equals, false)
		fmt.Printf("%#U\n", ch)
	}
	fmt.Println("'")

	// hiragana 3041 through 3093
	var hiragana []rune
	for r := rune(0x3041); r < rune(0x3094); r++ {
		fmt.Printf("%c", r)
		if r == 0x304a || r == 0x3054 || r == 0x305e || r == 0x3069 ||
			r == 0x306e || r == 0x307d || r == 0x3082 || r == 0x3088 ||
			r == 0x308d {

			fmt.Println()
		}
		hiragana = append(hiragana, r)
	}
	fmt.Println()
	strH := string(hiragana)
	// hiragana are 3 bytes each
	c.Assert(3*len(hiragana), Equals, len(strH))

	// hatakana 30a1 through 30ef
	var katakana []rune
	for r := rune(0x30a1); r < 0x30f0; r++ {
		fmt.Printf("%c", r)
		if r == 0x30aa || r == 0x30b4 || r == 0x30be || r == 0x30c9 ||
			r == 0x30ce || r == 0x30d4 || r == 0x30dc || r == 0x30e2 ||
			r == 0x30e8 || r == 0x30ed {

			fmt.Println()
		}
		katakana = append(katakana, r)
	}
	fmt.Println()
	strK := string(katakana)
	c.Assert(3*len(katakana), Equals, len(strK))
}
