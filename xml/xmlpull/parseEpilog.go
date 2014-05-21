package xmlpull

// xmlpull/parseEpilog.go

import (
	"fmt"
	"io"
	u "unicode"
)

var _ = fmt.Print

// The 'epilog' is the optional Misc* at the end of the prolog production.
//   22 prolog ::= XMLDecl? Misc* (doctypedecl Misc*)?
//
// XXX The existing code expects this to return just an error, but it modifies
// p.curEvent.
func (p *Parser) parseEpilog() (err error) {

	var (
		ch       rune
		curEvent PullEvent
		gotS     bool
	)

	if p.curEvent == END_DOCUMENT {
		err = p.NewXmlPullError("already reached end of XML input")
	} else {
		if p.reachedEnd {
			curEvent = END_DOCUMENT
		} else {
			gotS = false
			normalizeIgnorableWS := p.tokenizing && !p.roundtripSupported
			normalizedCR := false

			// epilog: Misc*
			ch, err = p.NextCh()
			p.afterLT = false // ???
			if !p.reachedEnd {
				for err != nil {
					// deal with Misc
					// [27] Misc ::= Comment | PI | S
					if ch == '<' {
						if gotS && p.tokenizing {
							p.afterLT = true
							curEvent = IGNORABLE_WHITESPACE
							break
						}
						ch, err = p.NextCh()
						if p.reachedEnd { // ????
							break
						}
						if ch == '?' {
							// check if it is 'xml'
							// deal with XMLDecl
							p.parsePI()
							if p.tokenizing {
								curEvent = PROCESSING_INSTRUCTION
								break
							}

						} else if ch == '!' {
							ch, err = p.NextCh()
							if err != nil || p.reachedEnd {
								break
							}
							if ch == 'D' {
								p.parseDocdecl() //FIXME
								if p.tokenizing {
									curEvent = DOCDECL
									break
								}
							} else if ch == '-' {
								p.parseComment()
								if p.tokenizing {
									curEvent = COMMENT
									break
								}
							} else {
								err = p.NewXmlPullError(
									"unexpected markup <!" + printableChar(ch))
							}
						} else if ch == '/' {
							err = p.NewXmlPullError(
								"end tag not allowed in epilog but got " + printableChar(ch))
						} else if isNameStartChar(ch) {
							err = p.NewXmlPullError(
								"start tag not allowed in epilog but got " + printableChar(ch))
						} else {
							err = p.NewXmlPullError(
								"in epilog expected ignorable content and not " + printableChar(ch))
						}
					} else if u.IsSpace(ch) {
						gotS = true
						if normalizeIgnorableWS {
							if ch == '\r' {
								normalizedCR = true
							} else if ch == '\n' {
								normalizedCR = false
							} else {
								normalizedCR = false
							}
						}
					} else {
						err = p.NewXmlPullError(
							"in epilog non whitespace content is not allowed but got " + printableChar(ch))
					}
					ch, err = p.NextCh()
					if p.reachedEnd {
						break
					}

				}
			}
			_ = normalizedCR // NOT USED XXX

			// throw Exception("unexpected content in epilog
		}
		if err == io.EOF {
			err = nil
			p.reachedEnd = true
		}

		if p.reachedEnd {
			if p.tokenizing && gotS {
				curEvent = IGNORABLE_WHITESPACE
			} else {
				curEvent = END_DOCUMENT
			}
		} else {
			err = p.NewXmlPullError("internal error in parseEpilog")
		}
	}
	if err == nil {
		p.curEvent = curEvent
	}
	return
}
