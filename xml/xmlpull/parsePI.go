package xmlpull

import (
	"fmt"
	"io"
	//"strings"
)

var _ = fmt.Print

// XML 1.0 Section 2.6 Processing Instructions
//
// [16] PI ::= '<?' PITarget (S (Char* - (Char* '?>' Char*)))? '?>'
// [17] PITarget         ::=    Name - (('X' | 'x') ('M' | 'm') ('L' | 'l'))
//
// Enter having seen '<?'.  Expect to see a target language name, whitespace,
// and then the body of the PI followed by '?>'.
//
// Return true if PI seen, false if an XmlDecl seen, error otherwise.  If
// true, p.piTarget contains the language name and b.piChars contains the
// body of the PI.
//
func (p *Parser) parsePI() (isPI bool, err error) {

	// Enter having seen  '<?'

	isPI = true // default assumption
	var (
		endSeen  bool
		piChars  []rune
		piTarget []rune
	)
	ch, err := p.NextCh()
	if err == nil {
		p.start()

		if p.IsS(ch) {
			err = p.NewXmlPullError("PI target may not begin with white space")
		}
		if err == nil {
			// collect the target language
			for err == nil {
				piTarget = append(piTarget, ch)
				ch, err = p.NextCh()
				if err == nil {
					if p.IsS(ch) {
						p.SkipS()
						break
					}
				}
			}
		}
		// We have a possibly nonsensical target.
		if err == nil && len(piTarget) == 3 {
			// Check for "xml" in any combo of cases but then
			// regard any capitalization as an error.
			if (piTarget[0] == 'x' || piTarget[0] == 'X') &&
				(piTarget[1] == 'm' || piTarget[1] == 'M') &&
				(piTarget[2] == 'l' || piTarget[2] == 'L') {

				// Is an XmlDecl and so must be must be right
				// at the start of the first line.  We have seen
				// "<?xml "
				if p.LineNo() != 1 || p.ColNo() != 6 {
					err = XmlDeclMustBeAtStart
				} else if string(piTarget) != "xml" {
					err = XmlInDeclMustBeLowerCase
				} else {
					err = p.parseXmlDecl()
					if err == nil {
						isPI = false
						endSeen = true
					}
					// XXX have not collected xmlDeclContent
				}
			}
		}

		// Unless it's an XmlDecl, collect the body.
		if isPI && err == nil {
			var haveQMark bool
			ch, err = p.NextCh()
			for err == nil {
				if ch == '?' {
					haveQMark = true
				} else if ch == '>' {
					if haveQMark {
						endSeen = true
						break
					} else {
						piChars = append(piChars, '>')
					}
					haveQMark = false

				} else {
					if haveQMark {
						piChars = append(piChars, '?')
						haveQMark = false
					}
					piChars = append(piChars, ch)
				}
				ch, err = p.NextCh()
			}
		}
	}
	if (err == nil && !endSeen) ||
		// not a terribly robust test
		// (err != nil && strings.HasSuffix(err.Error(), ": EOF")) {
		err == io.EOF {

		err = p.NotClosedErr("processing instruction")
	}
	if err == nil {
		if isPI {
			p.piChars = make([]rune, len(piChars))
			copy(p.piChars, piChars)
			p.piTarget = make([]rune, len(piTarget))
			copy(p.piTarget, piTarget)
		}
		// XXX If XmlDecl, that routine responsible for collecting its runes.
	}
	return
}
