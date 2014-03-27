package xmlpull

import (
	"fmt"
	// gu "github.com/jddixon/xgo/util"		// for INTERNING
)

var _ = fmt.Print

// entity reference http://www.w3.org/TR/2000/REC-xml-20001006#NT-Reference
// [67] Reference          ::=          EntityRef | CharRef
//
func (p *Parser) parseEntityRef() (runes []rune, err error) {

	// Enter having seen '&'

	ch, err := p.NextCh()
	if err != nil {
		return
	}

	if ch == '#' {
		// CHARACTER REFERENCE --------------------------------------
		var charRef rune
		ch, err = p.NextCh()
		if err == nil {
			return
		}
		if ch == 'x' {
			//encoded in hex
			for err == nil {
				ch, err = p.NextCh()

				if ch >= '0' && ch <= '9' {
					charRef = charRef*16 + (ch - '0')
				} else if ch >= 'a' && ch <= 'f' {
					charRef = charRef*16 + (ch - ('a' - 10))
				} else if ch >= 'A' && ch <= 'F' {
					charRef = charRef*16 + (ch - ('A' - 10))
				} else if ch == ';' {
					break
				} else {
					err = p.NewXmlPullError(
						"character reference (with hex value) may not contain " +
							printableChar(ch))
				}
			}
		} else {
			// encoded in decimal
			for err == nil {
				if ch >= '0' && ch <= '9' {
					charRef = charRef*10 + (ch - '0')
				} else if ch == ';' {
					break
				} else {
					err = p.NewXmlPullError(
						"character reference (with decimal value) may not contain " + printableChar(ch))
				}
				if err == nil {
					ch, err = p.NextCh()
				}
			}
		}
		runes := []rune{charRef}
		if p.tokenizing {
			// p.text = gu.Intern(runes, 0, 1)
			p.text = runes
		}
	} else {
		// ENTITY REF -----------------------------------------------
		//
		// [68]     EntityRef          ::=          '&' Name ';'
		// scan name up to semicolon
		if !isNameStartChar(ch) {
			err = p.NewXmlPullError(
				"entity reference names can not start with character '" + printableChar(ch) + "'")
		}
		var rName []rune
		for err == nil {
			ch, err = p.NextCh()
			// XXX
			if ch == ';' {
				break
			}
			if !isNameChar(ch) {
				err = p.NewXmlPullError(
					"entity reference name can not contain character " + printableChar(ch) + "'")
			}
			rName = append(rName, ch)
		}

		// XXX HANDLE ANY ERROR

		// determine what rName maps to
		length := len(rName)

		if length == 2 && rName[0] == 'l' && rName[1] == 't' {
			if p.tokenizing {
				p.text = []rune{'<'}
			}
			//charRefOneCharBuf := []rune{'<'}
			// return charRefOneCharBuf
			runes = []rune{'<'}

		} else if length == 3 && rName[0] == 'a' &&
			rName[1] == 'm' &&
			rName[2] == 'p' {

			if p.tokenizing {
				p.text = []rune{'&'}
			}
			runes = []rune{'&'}
		} else if length == 2 &&
			rName[0] == 'g' &&
			rName[1] == 't' {

			if p.tokenizing {
				p.text = []rune{'>'}
			}
			runes = []rune{'>'}
		} else if length == 4 && rName[0] == 'a' &&
			rName[1] == 'p' &&
			rName[2] == 'o' &&
			rName[3] == 's' {

			if p.tokenizing {
				p.text = []rune{'\''}
			}
			runes = []rune{'\''}
		} else if length == 4 && rName[0] == 'q' &&
			rName[1] == 'u' &&
			rName[2] == 'o' &&
			rName[3] == 't' {

			if p.tokenizing {
				p.text = []rune{'"'}
			}
			runes = []rune{'"'}
		} else {
			runes, err = p.lookupEntityReplacement(rName)
		}
		if err != nil || runes == nil || len(runes) == 0 {
			if p.tokenizing {
				p.text = nil
			}
		}
	}
	return
} // GEEP

func (p *Parser) lookupEntityReplacement(rName []rune) (
	replacement []rune, err error) {

	entityNameLen := uint(len(rName))

	// XXX KUKEMAL
	if false { // !p.allStringsInterned {
		hash := FastHash(rName)
		// rNameLen := uint(len(rName))
	DIJKSTRA:
		for i := p.entityEnd - 1; i >= 0; i-- {
			if hash == p.entityNameHash[i] &&
				entityNameLen == uint(len(p.entityNameBuf[i])) {

				entityBuf := p.entityNameBuf[i]
				for j := uint(0); j < entityNameLen; j++ {
					if rName[j] != entityBuf[j] {
						goto DIJKSTRA // no match
					}
				}
				if p.tokenizing {
					p.text = make([]rune, len(p.entityReplacement[i]))
					copy(p.text, p.entityReplacement[i])
				}
				replacement = make([]rune, len(p.entityReplacement[i]))
				copy(replacement, p.entityReplacement[i])
				return
			}
		}
	} else {
		//p.entityRefName = gu.Intern(rName)
		p.entityRefName = rName
		for i := p.entityEnd - 1; i >= 0; i-- {
			if SameRunes(p.entityRefName, p.entityName[i]) {
				if p.tokenizing {
					p.text = make([]rune, len(p.entityReplacement[i]))
					copy(p.text, p.entityReplacement[i])
				}
				replacement = make([]rune, len(p.entityReplacement[i]))
				copy(replacement, p.entityReplacement[i])
				return
			}
		}
	}
	return
}
