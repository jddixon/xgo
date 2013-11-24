package lex

// xgo/lex/lex_input_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
	"unicode/utf8"
)
var _ = fmt.Print

func (s *XLSuite) TestS (c *C) {

	whitespace := []byte(" \n\r\t")
	for i := 0; i < len(whitespace) ; i++ {
		c.Assert(IsS(whitespace[i]), Equals, true)
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
		fmt.Printf("%#U\n", ch)
	}
	fmt.Println("'")
}
