package xmlpull

import (
	"fmt"
)

var _ = fmt.Print

// [28]  doctypedecl ::= '<!DOCTYPE' S Name (S ExternalID)? S? ('['
//                      (markupdecl | DeclSep)* ']' S?)? '>'
//
func (p *Parser) parseDocdecl() (err error) {

	//ASSUMPTION: we have seen <!D

	err2 := p.ExpectStr("OCTYPE")
	if err2 != nil {
		err = p.NewXmlPullError(err2.Error())
	}

	// do simple and crude scanning for end of doctype

	bracketLevel := 0
	normalizeIgnorableWS := p.tokenizing && !p.roundtripSupported
	normalizedCR := false

	for err == nil {
		var ch rune
		ch, err = p.NextCh()
		if err != nil {
			break
		}
		if ch == '[' {
			bracketLevel++
		} else if ch == ']' {
			bracketLevel--
		}
		if ch == '>' && bracketLevel == 0 {
			break
		}
		if normalizeIgnorableWS {
			if ch == '\r' {
				normalizedCR = true
				//if !usePC {
				//    usePC = true
				//}
				//if pcEnd >= pc.length {
				//	ensurePC(pcEnd)
				// }
				// XXX NEED TO REPLACE WITH \n
				// pc[pcEnd++] = '\n'
			} else {
				normalizedCR = false
			}
		}

	}
	_ = normalizedCR // XXX not used
	return
}
