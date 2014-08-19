package xmlpull

// xgo/xml/xmlpull/next.go

import (
	"fmt"
	"io"
	"os"
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
		// DEBUG
		fmt.Printf("\nDO_NEXT: state == %s\n", PARSER_STATE_NAMES[p.state])
		// END
		switch p.state {
		case PRE_START_DOC:
			p.state = START_STATE
			curEvent = START_DOCUMENT
			return

		case START_STATE:
			// handle for xmlDecl is <?xml
			found, err = p.AcceptStr("<?xml")
			// DEBUG
			fmt.Printf("    <?xml found = %v\n", found)
			// END
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
			// DEBUG
			fmt.Println("fell through to check for XML_DECL_SEEN handles")
			fmt.Println("  checking for Misc1")
			// END

			// Misc1: ===============================================
			miscFound := true
			for err == nil && miscFound {
				miscFound, curEvent, err = p.acceptMisc()
				if miscFound && p.tokenizing {
					// DEBUG
					fmt.Printf("    MISC FOUND; curEvent := %s\n",
						PULL_EVENT_NAMES[curEvent])
					// END
					return
				}
			}
			if err != nil {
				return
			}

			// END Misc1 ============================================

			// DEBUG
			fmt.Println("  checking for DOCDECL")
			// END

			// handle for doctypedecl is '<!D
			var docDeclFound bool
			docDeclFound, err = p.AcceptStr("<!D")
			if err == nil {
				if docDeclFound {
					// DEBUG
					fmt.Println("docDeclFound DOC_TYPE_DECL")
					// END
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
			if !docDeclFound {
				p.state = EXPECT_START_ROOT
				continue
			}

			// DEBUG
			fmt.Println("falling through to DOC_DECL_SEEN")
			// END

			fallthrough

		case DOC_DECL_SEEN:

			// Misc2: ===============================================
			miscFound := true
			for err == nil && miscFound {
				miscFound, curEvent, err = p.acceptMisc()
				if miscFound && p.tokenizing {
					return
				}
			}
			if err != nil {
				return
			}
			p.state = EXPECT_START_ROOT
			fallthrough

		case EXPECT_START_ROOT:
			// DEBUG
			fmt.Println("EXPECT_START_ROOT: looking for START_ROOT handles")
			// END
			ch, err = p.NextCh()
			if err == nil {
				// handle for rootStart is '<'
				if ch == '<' {
					curEvent, err = p.parseStartTag()
					if (err == nil || err == io.EOF) && curEvent == START_TAG {
						if p.isEmptyElement {
							p.state = COLLECTING_EPILOG
						} else {
							p.state = START_ROOT_SEEN
						}
					} else {
						// UNHANDLED ERROR; should go to error state
					}
				} else {
					err = MissingRootElement
					p.state = ERROR_STATE
				}
			}
			return

		case START_ROOT_SEEN:
			// DEBUG
			fmt.Println("looking for START_ROOT_SEEN handles")
			// END

			// XXX STUB XXX

			// deeper handlers

			// XXX STUB XXX THIS IS WHERE WE PROCESS THE BODY OF THE
			// XXX ELEMENT

			// handle for rootEnd is '/>'
			ch, err = p.NextCh()
			fmt.Printf("ROOT END: ch '%c' err %s\n", ch, err) // DEBUG
			if err == nil && ch == '/' {
				ch, err = p.NextCh()
				fmt.Printf("ch '%c' err %s\n", ch, err) // DEBUG
				if (err == nil || err == io.EOF) && ch == '>' {
					fmt.Println("GOT >") //DEBUG
					curEvent, err = p.parseEndTag()
					// DEBUG
					fmt.Printf("parseEndTag returns curEvent %d, err %s\n",
						curEvent, err.Error)
					// END
				}
			}
			fmt.Printf("   START_ROOT_SEEN: state --> END_ROOT_SEEN\n")
			if p.isEmptyElement {
				p.state = COLLECTING_EPILOG
			} else {
				p.state = END_ROOT_SEEN
			}
			return

		case END_ROOT_SEEN:
			// accept zero or more Misc and EOF
			curEvent, err = p.parseEpilog()
			if err == io.EOF {
				p.state = PAST_END_DOC
			} else {
				p.state = COLLECTING_EPILOG
			}
			return

		case COLLECTING_EPILOG:
			// accept zero or more Misc and EOF
			curEvent, err = p.parseEpilog()
			if err == io.EOF {
				p.state = PAST_END_DOC
			} else {
				p.state = COLLECTING_EPILOG
			}
			return

		case ERROR_STATE:
			fmt.Printf("IRRECOVERABLE ERROR %s", err.Error)
			os.Exit(1)
		} // end switch
	} // end for

	return
}
