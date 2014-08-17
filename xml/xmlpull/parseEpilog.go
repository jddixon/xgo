package xmlpull

// xmlpull/parseEpilog.go

import (
	"fmt"
	"io"
)

var _ = fmt.Print

// The 'epilog' is the optional Misc* at the end of the prolog production.
//   22 prolog ::= XMLDecl? Misc* (doctypedecl Misc*)?
//
// XXX The existing code expects this to return just an error, but it modifies
// p.curEvent.
func (p *Parser) parseEpilog() (curEvent PullEvent, err error) {

	var (
		ch           rune
		gotS         bool // NOT USED
		normalizedCR bool // NOT USED
	)

	gotS = false
	normalizeIgnorableWS := p.tokenizing && !p.roundtripSupported

	// epilog: Misc*
	ch, err = p.NextCh()
	for err == nil || err == io.EOF {
		// deal with Misc
		// [27] Misc ::= Comment | PI | S
		if ch == '<' {
			ch, err = p.NextCh()
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
				if err != nil && err != io.EOF {
					break
				}
				if ch == 'D' {
					err = p.parseDocTypeDecl() //FIXME
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
		} else if p.IsS(ch) {
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
	}
	if err == nil || err == io.EOF {
		p.curEvent = curEvent
		if err != io.EOF {
			p.state = COLLECTING_EPILOG
		} else {
			p.state = PAST_END_DOC
		}
	}
	_, _ = gotS, normalizedCR // UNUSED
	return
}
