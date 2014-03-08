package xmlpull

import (
	// e "errors"
	"fmt"
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
	if p.afterLT {
		ch = buf[pos-1]
	} else {
		ch, err = p.NextCh()
	}

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

	p.afterLT = false

	for err == nil {
		// deal with Misc
		// deal with docdecl --> mark it!
		// else parseStartTag seen <[^/]
		if ch == '<' {
			if gotS && p.tokenizing {
				posEnd = pos - 1
				p.afterLT = true
				p.curEvent = IGNORABLE_WHITESPACE
				return //XXX  DIJKSTRA XXX
			}
			ch = more()
			if ch == '?' {
				// check if it is 'xml'
				// deal with XMLDecl
				if parsePI() { // make sure to skip XMLDecl
					if p.tokenizing {
						p.curEvent = PROCESSING_INSTRUCTION
						return // XXX DIJKSTRA
					}
				} else {
					// skip over - continue tokenizing
					posStart = pos
					gotS = false
				}

			} else if ch == '!' {
				ch = more()
				if ch == 'D' {
					if seenDocdecl {
						err = p.NewXmlPullError(
							"only one docdecl allowed in XML document", this, null)
						return // XXX DIJKSTRA
					}
					seenDocdecl = true
					parseDocdecl()
					if p.tokenizing {
						p.curEvent = DOCDECL
						return // DIJKSTRA
					}
				} else if ch == '-' {
					parseComment()
					if p.tokenizing {
						p.curEvent = COMMENT
						return // XXX DIJKSTRA
					}
				} else {
					err = p.NewXmlPullError(
						"unexpected markup <!"+printable(ch), this, null)
					return // DIJKSTRA
				}
			} else if ch == '/' {
				err = p.NewXmlPullError(
					"expected start tag name and not "+printable(ch), this, null)
				return // DIJKSTRA
			} else if isNameStartChar(ch) {
				p.haveRootTag = true
				parseStartTag() // XXX DIJKSTRA
				return
			} else {
				err = p.NewXmlPullError(
					"expected start tag name and not "+printable(ch), this, null)
				return // DIJKSTRA
			}
		} else if isS(ch) {
			gotS = true
			if normalizeIgnorableWS {
				if ch == '\r' {
					normalizedCR = true
					//posEnd = pos -1
					//joinPC()
					// posEnd is already is set
					if !usePC {
						posEnd = pos - 1
						if posEnd > posStart {
							joinPC()
						} else {
							usePC = true
							pcStart = 0
							pcEnd = 0
						}
					}
					//assert usePC == true
					if pcEnd >= pc.length {
						ensurePC(pcEnd)
					}
					pc[pcEnd] = '\n'
					pcEnd++
				} else if ch == '\n' {
					if !normalizedCR && usePC {
						if pcEnd >= pc.length {
							ensurePC(pcEnd)
						}
						pc[pcEnd] = '\n'
						pcEnd++
					}
					normalizedCR = false
				} else {
					if usePC {
						if pcEnd >= pc.length {
							ensurePC(pcEnd)
						}
						pc[pcEnd] = ch
						pcEnd++
					}
					normalizedCR = false
				}
			}
		} else {
			err = p.NewXmlPullError(
				"only whitespace content allowed before start tag and not "+printable(ch),
				this, null)
			return // DIJKSTRA
		}
		ch, err = p.NextCh()
	}
	return
}
