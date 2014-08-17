package xmlpull

import (
	"fmt"
	"io"
)

var _ = fmt.Print

// [28]  doctypedecl ::= '<!DOCTYPE' S Name (S ExternalID)? S? ('['
//                      (markupdecl | DeclSep)* ']' S?)? '>'
//
func (p *Parser) parseDocTypeDecl() (err error) {

	//ASSUMPTION: we have seen <!D

	err = p.ExpectStr("OCTYPE")
	if err != nil {
		err = p.NewXmlPullError(err.Error())
	}
	if err == nil {
		err = p.ExpectS()
	}

	// do simple and crude scanning for end of doctype
	var decl []rune
	if err == nil {
		bracketLevel := 0
		// normalizeIgnorableWS := p.tokenizing && !p.roundtripSupported
		// normalizedCR := false

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
			//if normalizeIgnorableWS {
			//	...
			// }
			decl = append(decl, ch)
		}
	}
	if err == nil || err == io.EOF {
		p.docTypeDecl = string(decl)
	}
	return
}
