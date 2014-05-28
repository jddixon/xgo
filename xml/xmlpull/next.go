package xmlpull

// xgo/xml/xmlpull/next.go

import (
	"fmt"
	// "strings"
)

var _ = fmt.Print

func (p *Parser) Next() (PullEvent, error) {
	p.tokenizing = false
	return p.doNext()
}

func (p *Parser) NextToken() (PullEvent, error) {
	p.tokenizing = true
	return p.doNext()
}

// Parser.state is initialized in NewParser()
func (p *Parser) doNext() (curEvent PullEvent, err error) {

	var (
		ch    rune
		found bool
	)

	for err == nil {
		p.text = p.text[:0] // clear the slice
		switch p.state {
		case PRE_START_DOC:
			p.state = START_STATE
			curEvent = START_DOCUMENT
			return

		case START_STATE:
			// handle for xmlDecl is <?xml
			found, err = p.AcceptStr("<?xml")
			if err == nil {
				if found {
					err = p.parseXmlDecl()
					if err == nil {
						p.state = XML_DECL_SEEN
					}
				}
			}
			if err != nil {
				return
			}
			fallthrough

		case XML_DECL_SEEN:
			// misc1:
			// handle for comment is '<!-'
			found, err = p.AcceptStr("<!-")
			if err == nil {
				if found {
					err = p.parseComment()
					if err == nil {
						curEvent = COMMENT
						if p.tokenizing {
							return
						} else {
							// no change in state
							continue
						}
					}
				}
			}
			if err != nil {
				return
			}
			// handle for PI is '<?'
			found, err = p.AcceptStr("<?")
			if err == nil {
				if found {
					found, err = p.parsePI()
					if err == nil && found {
						curEvent = PROCESSING_INSTRUCTION
						if p.tokenizing {
							return
						} else {
							// no change in state
							continue
						}
					}
				}
			}
			if err != nil {
				return
			}
			// handle for S is IsS()
			ch, err = p.PeekCh()
			for err == nil && p.IsS(ch) {
				p.text = append(p.text, ch) // ACCUMULATING WHITESPACE IN text
			}
			if err != nil {
				return
			}
			// POSSIBLE EOF?
			p.PushBack(ch)
			if len(p.text) > 0 {
				curEvent = IGNORABLE_WHITESPACE
				if p.tokenizing {
					return
				} else {
					// no change in state
					continue
				}
			}
			// handle for doctypedecl is '<!D
			found, err = p.AcceptStr("<!D")
			if err == nil {
				if found {
					err = p.parseDocTypeDecl()
					if err == nil {
						curEvent = DOCDECL
						// change in state
						p.state = DOC_DECL_SEEN
						if p.tokenizing {
							return
						} else {
							continue
						}
					}
				}
			}
			if err != nil {
				return
			}

			fallthrough

		case DOC_DECL_SEEN:

			// misc2: handles for comment, PI, S
			// handle for rootStart is '<' plus start char

		case START_ROOT_SEEN:
			// deeper handlers

			// handle for rootEnd is '/>'

		case END_ROOT_SEEN:
			// miscN: handlers for comment, PI, S
			// otherwise require EOF
			curEvent = END_DOCUMENT
			p.state = PAST_END_DOC

		case PAST_END_DOC:

		}
	}

	// VESTIGIAL CODE ===============================================
	if p.pastEndTag {
		p.pastEndTag = false
		p.elmDepth--

		// fewer namespaces available
		p.namespaceEnd = p.elNamespaceCount[p.elmDepth]
	}

	// XXX else ?
	if p.isEmptyElement { // set in parseStartTag
		p.isEmptyElement = false
		p.pastEndTag = true
		return END_TAG, nil // RETURN END_TAG
	}

	if p.elmDepth == 0 {
		if !p.rootElmSeen {
			// we haven't seen the root element yet
			err = p.parseProlog()
		} else {
			// we are post the root element; almost done
			err = p.parseEpilog()
		}
	} else { //  p.elmDepth > 0 {

		if p.seenStartTag {
			p.seenStartTag = false
			return p.parseStartTag() // RETURN START_TAG OR ERROR
		}
		if p.seenEndTag {
			p.seenEndTag = false
			return p.parseEndTag() // RETURN END_TAG OR ERROR
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

		// When true there is some potential event TEXT to return - keep
		// gathering
		charDataSeen := false

		// when true TEXT data is not continual (like <![CDATA[text]]>) and
		// requires PC merging

		var needsMerging bool

		for err == nil {
			// work on MARKUP
			if ch == '<' {
				if charDataSeen {
					if p.tokenizing {
						p.afterLT = true
						curEvent = TEXT
						p.curEvent = curEvent
						return // RETURN TEXT
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
					return // RETURN TEXT OR END_TAG OR ERROR
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
							return // COMMENT
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
							return // RETURN CDSECT
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
						return // RETURN PROCESSING_INSTRUCTION
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
						return // RETURN TEXT
					}
					curEvent, err = p.parseStartTag()
					if err == nil {
						p.curEvent = curEvent
					}
					return // RETURN START_TAG OR ERROR
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
					return // RETURN TEXT
				}
				var resolvedEntity []rune
				resolvedEntity, err = p.parseEntityRef()
				if err != nil {
					break
				}
				if p.tokenizing {
					curEvent = ENTITY_REF
					p.curEvent = curEvent
					return // RETURN ENTITY_REF
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
	}
	return
}
