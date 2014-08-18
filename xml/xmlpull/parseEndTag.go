package xmlpull

import (
	"fmt"
	"io"
)

var _ = fmt.Print

// [42] ETag ::=  '</' Name S? '>'

func (p *Parser) parseEndTag() (curEvent PullEvent, err error) {

	// ASSUMPTION input is past "</"

	fmt.Println("Entering parseEndTag") // DEBUG

	ch, err := p.NextCh()
	if !isNameStartChar(ch) {
		msg := fmt.Sprintf(
			"expected name start character, found '%c", ch)
		err = p.NewXmlPullError(msg)
	}
	if err == nil {
		var name []rune
		name = append(name, ch)
		for err == nil {
			ch, err = p.NextCh()
			if err == nil && isNameChar(ch) {
				name = append(name, ch)
			} else {
				break
			}
		}
		// DEBUG
		fmt.Printf("parseEndTag: name is %s\n", name)
		// END
		if err == nil {
			// end tag must match start tag
			startName := p.elRawName[p.elmDepth]
			if !SameRunes(name, startName) {
				msg := fmt.Sprintf("end tag %s does not match start tag %s\n",
					string(name), string(startName))
				err = p.NewXmlPullError(msg)
			}
			if err == nil {
				fmt.Printf("COLLECTING /> FOR ROOT TAG\n") // DEBUG
				p.SkipS()
				ch, err = p.NextCh()
				if err == nil && ch != '/' {
					msg := fmt.Sprintf(
						"expected '/' to finish end tag not '%c'", ch)
					err = p.NewXmlPullError(msg)
				}
				if err == nil {
					ch, err = p.NextCh()
					if (err == nil || err == io.EOF) && ch != '>' {
						msg := fmt.Sprintf(
							"expected '>' to finish end tag not '%c'", ch)
						err = p.NewXmlPullError(msg)
					}
				}
				if err != nil || err == io.EOF {
					curEvent = END_TAG
				}
			}
		}
	}
	return
}
