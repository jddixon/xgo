package xmlpull

// xgo/xml/xmlpull/parseProlog.go

import (
	// e "errors"
	"fmt"
	u "unicode"
)

var _ = fmt.Print

// [1]  document ::= prolog element Misc*
// [2]  prolog: ::= XMLDecl? Misc* (doctypedecl Misc*)?
// [27] Misc ::= Comment | PI | S
// [28] doctypedecl ::= '<!DOCTYPE' S Name (S ExternalID)? S? ('['
//                      (markupdecl | DeclSep)* ']' S?)? '>'
// [39] element ::= EmptyElemTag | STag content ETag
// [40] STag ::= '<' Name (S Attribute)* S? '>'
//
// Enter having already parsed any XmlDecl.
//
func (p *Parser) parseRestOfProlog() (err error) {

	var ch rune
	//if p.afterLT {
	//	ch = buf[pos-1]
	// } else {
	ch, err = p.NextCh()
	//}

	// DOES NOT BELONG HERE. ----------------------------------------
	// This block analyzes the very first byte in a document.
	if err == nil && p.curEvent == START_DOCUMENT {
		// This is the first character of input, and so might be the
		// unicode byte order mark (BOM)
		if ch == '\uFFFE' {
			panic("data in wrong byte order!")
		} else if ch == '\uFEFF' {
			// discard
			ch, err = p.NextCh()
		}
	}
	// END DOES NOT BELONG ------------------------------------------
	if err == nil {
		var normalizedCR bool // should this be here ?? XXX
		p.afterLT = false
		gotS := false

		for err == nil {
			// deal with Misc
			// deal with docdecl --> mark it!
			// else parseStartTag seen <[^/]
			if ch == '<' {
				if gotS && p.tokenizing {
					// posEnd = pos - 1
					p.afterLT = true
					p.curEvent = IGNORABLE_WHITESPACE
					break
				}
				ch, err = p.NextCh()
				if ch == '?' {
					// check if it is 'xml'
					// deal with XMLDecl
					var isPI bool
					isPI, err = p.parsePI() // skipping XMLDecl
					if err != nil {
						break
					}
					if isPI {
						if p.tokenizing {
							p.curEvent = PROCESSING_INSTRUCTION
							break
						}
					} else {
						// skip over - continue tokenizing
						//posStart = pos
						gotS = false
					}

				} else if ch == '!' {
					ch, err = p.NextCh()
					if ch == 'D' {
						if p.seenDocdecl {
							err = p.NewXmlPullError(
								"only one docdecl allowed in XML document")
							break
						}
						p.seenDocdecl = true
						p.parseDocdecl()
						if p.tokenizing {
							p.curEvent = DOCDECL
							break
						}
					} else if ch == '-' {
						p.parseComment()
						if p.tokenizing {
							p.curEvent = COMMENT
							break
						}
					} else {
						err = p.NewXmlPullError(
							"unexpected markup <!" + printableChar(ch))
						break
					}
				} else if ch == '/' {
					err = p.NewXmlPullError(
						"expected start tag name and not " + printableChar(ch))
					break
				} else if isNameStartChar(ch) {
					p.haveRootTag = true
					p.parseStartTag(ch)
					break
				} else {
					err = p.NewXmlPullError(
						"expected start tag name and not " + printableChar(ch))
					break
				}
			} else if u.IsSpace(ch) {
				gotS = true
				if p.normalizeIgnorableWS {
					if ch == '\r' {
						normalizedCR = true
						//posEnd = pos -1
						//joinPC()
						// posEnd is already is set
						//if !usePC {
						//	posEnd = pos - 1
						//	if posEnd > posStart {
						//		joinPC()
						//	} else {
						//		usePC = true
						//		pcStart = 0
						//		pcEnd = 0
						//	}
						//}
						//assert usePC == true
						//if pcEnd >= pc.length {
						//		ensurePC(pcEnd)
						//
						//pc[pcEnd] = '\n'
						//pcEnd++
					} else if ch == '\n' {
						//if !normalizedCR && usePC {
						//	if pcEnd >= pc.length {
						//		ensurePC(pcEnd)
						//	}
						//	pc[pcEnd] = '\n'
						//	pcEnd++
						//}
						normalizedCR = false
					} else {
						//if usePC {
						//	if pcEnd >= pc.length {
						//		ensurePC(pcEnd)
						//	}
						//	pc[pcEnd] = ch
						//	pcEnd++
						//}
						normalizedCR = false
					}
				}
			} else {
				err = p.NewXmlPullError(
					"only whitespace content allowed before start tag and not " + printableChar(ch))
				break
			}
			ch, err = p.NextCh()
		} // end for
		_ = normalizedCR // XXx NEVER USED
	}

	return
}
