package xmlpull

// xgo/xml/xmlpull/parseXmlDecl.go

import (
	e "errors"
	"fmt"
)

// ------------------------------------------------------------------
// Beware: parseXmlDecl() is called by parseProcessingInstruction()
// because it is a specific type of PI.
// ------------------------------------------------------------------

// [25] Eq ::= S? '=' S?

func (p *Parser) expectEq() (err error) {
	p.SkipS()
	ch, err := p.NextCh()
	if err == nil && ch != '=' {
		msg := fmt.Sprintf(
			"parseXmlDecl.expectEq: expected = , found %c\n", ch)
		err = e.New(msg)
	}
	if err == nil {
		p.SkipS() // closes Eq
	}
	return
}
func (p *Parser) expectQuoteCh() (quoteCh rune, err error) {
	quoteCh, err = p.NextCh()
	if err == nil && quoteCh != '\'' && quoteCh != '"' {
		msg := fmt.Sprintf("expected quotation mark, found '%c'", quoteCh)
		err = e.New(msg)
	}
	return
}

// [81] EncName ::= [A-Za-z] ([A-Za-z0-9._] | '-')*
func (p *Parser) getEncodingStartCh() (ch rune, err error) {
	ch, err = p.NextCh()
	if err == nil {
		if !('a' <= ch && ch <= 'z') && !('A' <= ch && ch <= 'Z') {
			msg := fmt.Sprintf("cannot start encoding name: '%c'\n", ch)
			err = e.New(msg)
		}
	}
	return
}

func (p *Parser) getEncodingNameCh(quoteCh rune) (ch rune, err error) {
	ch, err = p.NextCh()
	if err == nil && ch != quoteCh {
		if !('a' <= ch && ch <= 'z') && !('A' <= ch && ch <= 'Z') &&
			!('0' <= ch && ch <= '9') && (ch != '.') && (ch != '_') &&
			(ch != '-') {
			msg := fmt.Sprintf("illegal character in encoding name: '%c'\n", ch)
			err = e.New(msg)
		}
	}
	return
}

// Function called after encountering <?xmlS at the beginning of the input,
// where S as usual represents a space.
//
func (p *Parser) parseXmlDecl() (err error) {

	var (
		ch, quoteCh rune
	)

	// [23] XMLDecl ::= '<?xml' VersionInfo EncodingDecl? SDDecl? S? '?>'
	// [24] VersionInfo ::= S 'version' Eq ("'" VersionNum "'" | '"' VersionNum '"')
	// We are on first S past <?xml
	p.SkipS()
	err = p.ExpectStr("version")

	if err == nil {
		err = p.expectEq()
		if err == nil {
			quoteCh, err = p.expectQuoteCh()
		}
	}

	// [26] VersionNum ::= ([a-zA-Z0-9_.:] | '-')+
	if err == nil {
		var vRunes []rune
		ch, err = p.NextCh()
		for err == nil && ch != quoteCh {
			if ('a' <= ch && ch <= 'z') ||
				('A' <= ch && ch <= 'Z') ||
				('0' <= ch && ch <= '9') ||
				(ch == '_') || (ch == '.') || (ch == ':') || (ch == '-') {
				vRunes = append(vRunes, ch)
			} else {
				msg := fmt.Sprintf(
					"Not an acceptable version character: '%c'", ch)
				err = e.New(msg)
				break
			}
			ch, err = p.NextCh()
		}
		if err == nil {
			// ch is guaranteed to be quoteCh
			p.xmlDeclVersion = string(vRunes)
			if p.xmlDeclVersion != "1.0" {
				err = OnlyVersion1_0
			}
		}
	}

	// [80] EncodingDecl ::= S 'encoding' Eq ('"' EncName '"' | "'" EncName "'" )
	if err == nil {
		var found bool
		p.SkipS()
		found, err = p.AcceptStr("encoding")
		if err == nil {
			if found {
				var eRunes []rune
				err = p.expectEq()
				if err == nil {
					var eStartCh, eNameCh rune
					quoteCh, err = p.expectQuoteCh()
					if err == nil {
						eStartCh, err = p.getEncodingStartCh()
						if err == nil {
							eRunes = append(eRunes, eStartCh)
							for err == nil && eNameCh != quoteCh {
								eNameCh, err = p.getEncodingNameCh(quoteCh)
								if err == nil {
									if eNameCh == quoteCh {
										break
									}
									eRunes = append(eRunes, eNameCh)
								}
							}
						}
					}
				}
				if err == nil {
					p.xmlDeclEncoding = string(eRunes)
				}
			}
		}
	}
	// [32] SDDecl ::= S 'standalone' Eq (("'" ('yes' | 'no') "'") | ('"' ('yes' | 'no') '"'))
	if err == nil {
		var found bool
		p.SkipS()
		found, err = p.AcceptStr("standalone")
		if err == nil && found {
			var foundYes, foundNo bool
			err = p.expectEq()
			if err == nil {
				quoteCh, err = p.expectQuoteCh()
				if err == nil {
					foundYes, err = p.AcceptStr("yes")
				}
				if err == nil && !foundYes {
					foundNo, err = p.AcceptStr("no")
				}
				if err == nil {
					if foundYes {
						p.xmlDeclStandalone = true
					} else if foundNo {
						p.xmlDeclStandalone = false
					} else {
						err = MustBeYesOrNo
					}
				}
				ch, err = p.NextCh()
				if err == nil && ch != quoteCh {
					err = MissingDeclClosingQuote
				}
			}
		}
	}
	// expecting ?>
	if err == nil {
		p.SkipS()
		err = p.ExpectStr("?>")
	}
	return
}
