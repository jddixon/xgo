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

// [23] XMLDecl ::= '<?xml' VersionInfo EncodingDecl? SDDecl? S? '?>'
//
// Function called after encountering <?xml at the beginning of the input.
// THIS IS A CHANGE: handle was <?xmlS
//
func (p *Parser) parseXmlDecl() (err error) {

	var (
		ch, quoteCh rune
	)
	// We must be on the first S past <?xml

	// [24] VersionInfo ::= S 'version' Eq ("'" VersionNum "'" |
	//                                      '"' VersionNum '"')
	err = p.ExpectS()
	if err == nil {
		err = p.ExpectStr("version")
		if err == nil {
			err = p.expectEq() // [25]
			if err == nil {
				quoteCh, err = p._expectQuoteCh()
			}
		}
	}

	// [26] VersionNum ::= ([a-zA-Z0-9_.:] | '-')+
	//
	// Expect a valid version number; unless err, write it to p.xmlDeclVersion.
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

	// [80] EncodingDecl ::= S 'encoding' Eq ('"' EncName '"' |
	//                                        "'" EncName "'" )
	// [81] EncName ::= [A-Za-z] ([A-Za-z0-9._] | '-')*
	//
	// This production ([80]) is optional.  If found and no error, write
	// the encoding name (XXX as a string) to p.xlmDeclEncoding.
	//
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
					quoteCh, err = p._expectQuoteCh()
					if err == nil {
						eStartCh, err = p._getEncodingStartCh()
						if err == nil {
							eRunes = append(eRunes, eStartCh)
							for err == nil && eNameCh != quoteCh {
								eNameCh, err = p._getEncodingNameCh(quoteCh)
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
	// [32] SDDecl ::= S 'standalone' Eq (("'" ('yes' | 'no') "'") |
	//                                    ('"' ('yes' | 'no') '"'))
	//
	// This production is optional.  If it is present and there is no error,
	// set p.xmlDeclStandaone accordingly.  Otherwise this field defaults
	// to false.
	//
	if err == nil {
		var found bool
		p.SkipS()
		found, err = p.AcceptStr("standalone")
		if err == nil && found {
			var foundYes, foundNo bool
			err = p.expectEq()
			if err == nil {
				quoteCh, err = p._expectQuoteCh()
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
	// optional S, required '?>'
	if err == nil {
		p.SkipS()
		err = p.ExpectStr("?>")
	}
	return
}

// UTILITIES ========================================================

func (p *Parser) _expectQuoteCh() (quoteCh rune, err error) {
	quoteCh, err = p.NextCh()
	if err == nil && quoteCh != '\'' && quoteCh != '"' {
		msg := fmt.Sprintf("expected quotation mark, found '%c'", quoteCh)
		err = e.New(msg)
	}
	return
}

// [81] EncName ::= [A-Za-z] ([A-Za-z0-9._] | '-')*
//
func (p *Parser) _getEncodingStartCh() (ch rune, err error) {
	ch, err = p.NextCh()
	if err == nil {
		if !('a' <= ch && ch <= 'z') && !('A' <= ch && ch <= 'Z') {
			msg := fmt.Sprintf("cannot start encoding name: '%c'\n", ch)
			err = e.New(msg)
		}
	}
	return
}

func (p *Parser) _getEncodingNameCh(quoteCh rune) (ch rune, err error) {
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
