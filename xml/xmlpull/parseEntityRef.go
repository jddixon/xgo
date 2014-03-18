package xmlpull

import (
	"fmt"
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
		if p.tokenizing {
			p.text = []rune{charRef} // a copy, of course
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
		var name []rune
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
			name = append(name, ch)
		}

		// XXX HANDLE ANY ERROR

		// determine what name maps to
		length := len(name)

		if length == 2 && name[0] == 'l' && name[1] == 't' {
			if p.tokenizing {
				p.text = []rune{'<'}
			}
			//charRefOneCharBuf := []rune{'<'}
			// return charRefOneCharBuf
			runes = []rune{'<'}

		} else if length == 3 && name[0] == 'a' &&
			name[1] == 'm' &&
			name[2] == 'p' {

			if p.tokenizing {
				p.text = []rune{'&'}
			}
			runes = []rune{'&'}
		} else if length == 2 &&
			name[0] == 'g' &&
			name[1] == 't' {

			if p.tokenizing {
				p.text = []rune{'>'}
			}
			runes = []rune{'>'}
		} else if length == 4 && name[0] == 'a' &&
			name[1] == 'p' &&
			name[2] == 'o' &&
			name[3] == 's' {

			if p.tokenizing {
				p.text = []rune{'\''}
			}
			runes = []rune{'\''}
		} else if length == 4 && name[0] == 'q' &&
			name[1] == 'u' &&
			name[2] == 'o' &&
			name[3] == 't' {

			if p.tokenizing {
				p.text = []rune{'"'}
			}
			runes = []rune{'"'}
		} else {
			runes, err = p.lookupEntityReplacement(name)
		}
		if err != nil || runes == nil || len(runes) == 0 {
			if p.tokenizing {
				p.text = nil
			}
		}
	}
	return
} // GEEP

func (p *Parser) lookupEntityReplacement(name []rune) (

	replacement []rune, err error) {
	entityNameLen := uint(len(name))
	hash := FastHash(name)
	for i := p.entityEnd - 1; i >= 0; i-- {
		notAMatch := false
		if hash == p.entityNameHash[i] &&
			entityNameLen == uint(len(p.entityNameBuf[i])) {

			entityBuf := p.entityNameBuf[i]
			for j := uint(0); j < entityNameLen; j++ {
				if name[j] != entityBuf[j] {
					notAMatch = true
					break
				}
			}
			if !notAMatch {
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
