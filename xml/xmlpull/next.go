package xmlpull

// xgo/xml/xmlpull/next.go

import (
	"fmt"
	// "strings"
)

var _ = fmt.Print

// [1] document ::= prolog element Misc*
//
func (p *Parser) Next() (PullEvent, error) {
	p.tokenizing = false
	return p.doNext()
}

func (p *Parser) NextToken() (PullEvent, error) {
	p.tokenizing = true
	return p.doNext()
}

func (p *Parser) doNext() (curEvent PullEvent, err error) {

	p.text = p.text[:0] // clear the slice

	if p.pastEndTag {
		p.pastEndTag = false
		p.elmDepth--
		p.namespaceEnd = p.elNamespaceCount[p.elmDepth] // less namespaces available
	}
	if p.isEmptyElement {
		p.isEmptyElement = false
		p.pastEndTag = true
		return END_TAG, nil
	}

	if p.elmDepth > 0 {

		if p.seenStartTag {
			p.seenStartTag = false
			return p.parseStartTag() // XXx NEEDS ch
		}
		if p.seenEndTag {
			p.seenEndTag = false
			return p.parseEndTag()
		}

		// ASSUMPTION: we are _on_ first character of content or markup

		// [43] content ::= CharData? ((element | Reference | CDSect | PI | Comment) CharData?)*

		var ch rune
		if p.afterLT { // we have read ahead ...
			p.afterLT = false
			ch = '<'
		} else if p.seenAmpersand {
			p.seenAmpersand = false
			ch = '&'
		} else {
			ch, err = p.NextCh()
		}

		// when true there is some potential event TEXT to return - keep gathering
		charDataSeen := false

		// when true TEXT data is not continual (like <![CDATA[text]]>) and requires PC merging
		var needsMerging bool

		for err == nil {
			// work on MARKUP
			if ch == '<' {
				if charDataSeen {
					if p.tokenizing {
						p.afterLT = true
						curEvent = TEXT
						p.curEvent = curEvent
						return
					}
				}
				ch, err = p.NextCh()
				if err != nil {
					break
				}
				if ch == '/' {
					if !p.tokenizing && charDataSeen {
						p.seenEndTag = true
						curEvent = TEXT
					} else {
						curEvent, err = p.parseEndTag()
					}
					p.curEvent = curEvent
					return
				} else if ch == '!' {
					ch, err = p.NextCh()
					if err != nil {
						break
					}
					if ch == '-' {
						p.parseComment()
						if p.tokenizing {
							curEvent = COMMENT
							p.curEvent = curEvent
							return
						}
						//if !usePC && charDataSeen  {
						//    needsMerging = true
						//} else {
						//    //completely ignore comment
						//}
					} else if ch == '[' {
						err = p.parseCDSect(charDataSeen)
						if err != nil {
							break
						}
						if p.tokenizing {
							curEvent = CDSECT
							p.curEvent = curEvent
							return
						}

						// if !usePC {
						// needsMerging = true
						// }
						// }

					} else {
						err = p.NewXmlPullError(
							"unexpected character in markup " +
								printableChar(ch))
					}
				} else if ch == '?' {
					var isPI bool
					isPI, err = p.parsePI()
					_ = isPI // XXX
					// XXX HANDLE ERROR
					if p.tokenizing {
						curEvent = PROCESSING_INSTRUCTION
						p.curEvent = curEvent
						return
					}
					// if !usePC && charDataSeen {
					//	needsMerging = true
					//} else {
					//	//completely ignore PI
					//}

				} else if isNameStartChar(ch) {
					if !p.tokenizing && charDataSeen {
						p.seenStartTag = true
						curEvent = TEXT
						p.curEvent = curEvent
						return
					}
					curEvent, err = p.parseStartTag()
					if err == nil {
						p.curEvent = curEvent
					}
					return
				} else {
					err = p.NewXmlPullError(
						"unexpected character in markup " +
							printableChar(ch))
				}
				// do content compaction if it makes sense

			} else if ch == '&' {
				// work on ENTITTY
				if p.tokenizing && charDataSeen {
					p.seenAmpersand = true
					curEvent = TEXT
					p.curEvent = curEvent
					return
				}
				var resolvedEntity []rune
				resolvedEntity, err = p.parseEntityRef()
				if err != nil {
					break
				}
				if p.tokenizing {
					curEvent = ENTITY_REF
					p.curEvent = curEvent
					return
				}
				// check if replacement text can be resolved !!!
				if len(resolvedEntity) == 0 {
					if len(p.entityRefName) == 0 {
						p.entityRefName = "UNKNOWN" // XXX
					}
					err = p.NewXmlPullError(
						"could not resolve entity named '" +
							printableStr(string(p.entityRefName)) + "'")
				}
				//if !usePC {
				//    if(charDataSeen) {
				//        joinPC(); // posEnd is already set correctly!!!
				//        needsMerging = false
				//    } else {
				//        usePC = true
				//        pcStart = pcEnd = 0
				//    }
				//} // FOO

				// write into PC replacement text - do merge for replacement text!!!!
				//for i := 0; i < len(resolvedEntity); i++ {
				//    if(pcEnd >= pc.length {
				//		ensurePC(pcEnd)
				//	}
				//    pc[pcEnd++] = resolvedEntity[ i ]
				//}
				charDataSeen = true
			} else {

				//if needsMerging {
				//    //assert usePC == false
				//    joinPC()  // posEnd is already set correctly!!!
				//    //posStart = pos  -  1
				//    needsMerging = false
				//}

				//no MARKUP not ENTITIES so work on character data ...

				// [14] CharData ::=   [^<&]* - ([^<&]* ']]>' [^<&]*)

				charDataSeen = true

				normalizedCR := false
				normalizeInput := !p.tokenizing || !p.roundtripSupported
				// use loop locality here!!!!
				seenBracket := false
				seenBracketBracket := false
				for { // do {

					// check that ]]> does not show in
					if ch == ']' {
						if seenBracket {
							seenBracketBracket = true
						} else {
							seenBracket = true
						}
					} else if seenBracketBracket && ch == '>' {
						err = p.NewXmlPullError(
							"characters ]]> are not allowed in content")
					} else {
						if seenBracket {
							seenBracket = false
							seenBracketBracket = false
						}
					}
					if normalizeInput {
						// deal with normalization issues ...
						if ch == '\r' {
							normalizedCR = true
							// MISSING: ADD REPLACMENT \n
						} else if ch == '\n' {
							normalizedCR = false
						} else {
							normalizedCR = false
						}
					}

					ch, err = p.NextCh()
					if ch == '<' || ch == '&' {
						break
					}

				}
				_ = needsMerging // XXX
				_ = normalizedCR // XXX
				continue         // ie, continue outer loop using this ch
			}
			ch, err = p.NextCh()
		} // endless for err == nil
	} else {
		if p.haveRootTag {
			err = p.parseEpilog()
		} else {
			err = p.parseRestOfProlog()
		}
	}
	return
}
