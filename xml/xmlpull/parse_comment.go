package xmlpull

import (
	"fmt"
	"io"
)

var _ = fmt.Print

// XML 1.0 Section 2.5 Comments
//
// [15] Comment :== '<!--' ((Char - '-') | ('-' (Char - '-')))* '-->'
//
// Grammer does not allow ---> at end.
//
func (p *Parser) parseComment() (err error) {

	// Enter having seen "<!-"

	var (
		commentChars     []rune
		endOfCommentSeen bool
	)
	lx := p.GetLexer()

	ch, err := lx.NextCh()
	if ch != '-' {
		err = p.NewXmlPullError("comment must start with <!--")
	}
	if err == nil {
		p.startLine = p.LineNo()
		p.startCol = p.ColNo()
		haveDash := false
		haveTwoDashes := false
		for err == nil {
			ch, err = lx.NextCh()
			if err == nil && haveTwoDashes && ch != '>' {
				err = p.NewXmlPullError("cannot have two dashes within comment")
				break
			}
			if ch == '-' {
				if !haveDash {
					haveDash = true
				} else {
					haveTwoDashes = true
					haveDash = false
				}
			} else if ch == '>' {
				if haveTwoDashes {
					endOfCommentSeen = true
					break // end of comment
				} else {
					haveTwoDashes = false
				}
				haveDash = false
			} else {
				if haveDash {
					commentChars = append(commentChars, '-')
					haveDash = false
				}
				// \r, \n handled the same way
				commentChars = append(commentChars, ch)
			}
		}
	}
	if (err == nil && !endOfCommentSeen) || err == io.EOF {
		err = p.NotClosedErr("comment")
	}
	if err == nil {
		p.commentChars = string(commentChars)
	}
	return
}
