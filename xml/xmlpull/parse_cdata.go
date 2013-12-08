package xmlpull

import (
	"strings"
)

// XML 1.0 Section 2.7 CDATA Sections
//
// [18] CDSect ::= CDStart CData CDEnd
// [19] CDStart ::=  '<![CDATA['
// [20] CData ::= (Char* - (Char* ']]>' Char*))
// [21] CDEnd ::= ']]>'
//
func (p *Parser) parseCDSect() (err error) {

	// Enter having seen <![

	var (
		cDataChars                   []rune
		ch                           rune
		endSeen                      bool
		haveBracket, haveTwoBrackets bool
	)
	err = p.ExpectStr("CDATA[")

	if err == nil {
		p.start() // set up line & col numbers

		// loop until we see ]]>
		for err == nil {
			ch, err = p.NextCh()
			if err == nil {
				if ch == ']' {
					if !haveBracket {
						haveBracket = true
					} else {
						haveTwoBrackets = true
						haveBracket = false
					}
				} else if ch == '>' {
					if haveTwoBrackets {
						endSeen = true
						break
					} else {
						haveTwoBrackets = false
					}
					haveBracket = false
				} else {
					if haveBracket {
						cDataChars = append(cDataChars, ']')
						haveBracket = false
					} else if haveTwoBrackets {
						cDataChars = append(cDataChars, ']')
						cDataChars = append(cDataChars, ']')
						haveTwoBrackets = false
					}
					cDataChars = append(cDataChars, ch)
				}
			}
		}
	}
	if (err == nil && !endSeen) ||
		// not a terribly robust test
		(err != nil && strings.HasSuffix(err.Error(), ": EOF")) {
		err = p.NotClosedErr("cData")
	}
	if err == nil {
		p.cDataChars = string(cDataChars)
	}
	return
}
