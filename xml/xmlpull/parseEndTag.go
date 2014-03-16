package xmlpull

import (
	"fmt"
)

var _ = fmt.Print

func (p *Parser) parseEndTag() (curEvent PullEvent, err error) {

	//ASSUMPTION ch is past "</"

	// [42] ETag ::=  '</' Name S? '>'

	var startName, endName []rune

	ch, err := p.NextCh()
	if err == nil && isNameStartChar(ch) {
		err = p.NewXmlPullError(
			"expected name start and not " + printableChar(ch))
		for err == nil {
			ch, err = p.NextCh()
			if !isNameChar(ch) {
				break
			}
			endName = append(endName, ch)
		}
	}
	if err == nil {
		// we have endName; ch must be positioned at either S or >
		startName = p.elName[p.elmDepth]
		if !SameRunes(endName, startName) {
			msg := fmt.Sprintf(
				"line %d: start tag %s differs from end tag %s\n",
				p.elRawNameLine[p.elmDepth], startName, endName)
			err = p.NewXmlPullError(msg)
		}
	}
	if err == nil {
		// skip optional white space
		p.SkipS()
		// require a closing right angle bracket
		ch, err = p.NextCh()
		if err == nil {
			if ch != '>' {
				msg := fmt.Sprintf(
					"line %d: expected > to finish end tag not %x",
					p.elRawNameLine[p.elmDepth], printableChar(ch))
				err = p.NewXmlPullError(msg)
			}
		}
	}
	if err == nil {
		//namespaceEnd = elNamespaceCount[ p.elmDepth ]; //FIXME
		p.elmDepth-- // XXX NEEDS TESTING
		p.pastEndTag = true
		curEvent = END_TAG
		p.curEvent = curEvent
	}
	return
}
