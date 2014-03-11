package xmlpull

import (
	"fmt"
	u "unicode"
)

var _ = fmt.Print

// [40] STag ::=  '<' Name (S Attribute)* S? '>'
// [44] EmptyElemTag ::= '<' Name (S Attribute)* S? '/>'

// XXx the first name character is the parameter
//
func (p *Parser) parseStartTag(ch rune) (curEvent PullEvent, err error) {

	//ASSUMPTION ch is past <T

	p.elmDepth++
	// posStart = pos - 2
	p.isEmptyElement = false
	attributeCount := 0

	// retrieve name
	// nameStart := pos - 1 + bufAbsoluteStart
	// int colonPos = -1
	var colonFound bool
	if ch == ':' && p.processNamespaces {
		err = p.NewXmlPullError(
			"when namespaces processing enabled colon cannot be at element name start")
	}
	for err == nil {
		ch, err = p.NextCh()
		if err != nil {
			break
		}
		if !isNameChar(ch) {
			break
		}
		if ch == ':' && p.processNamespaces {
			if colonFound {
				err = p.NewXmlPullError(
					"only one colon is allowed in name of element when namespaces are enabled")
			}
			// colonPos = pos - 1 + bufAbsoluteStart
			colonFound = true
		}
	}

	// retrieve name
	// ensureElementsCapacity()			// XXX MUST IMPLEMENT

	// elLen := (pos - 1) - (nameStart - bufAbsoluteStart)
	elLen := 0 // XXX MUST DETERMINE VALUE

	if p.elRawName[p.elmDepth] == nil || len(p.elRawName[p.elmDepth]) < elLen {
		p.elRawName[p.elmDepth] = make([]rune, 2*elLen) // new char[ 2 * elLen ]
	}
	// XXX FIX ME
	// System.arraycopy(buf, nameStart - bufAbsoluteStart, p.elRawName[ p.elmDepth ], 0, elLen)
	p.elRawNameEnd[p.elmDepth] = elLen
	p.elRawNameLine[p.elmDepth] = p.lineNo

	name := ""

	// work on prefixes and namespace URI
	prefix := ""
	if p.processNamespaces {
		if colonFound {
			// XXX FIX ME FIX ME FIX ME
			// p.elPrefix[ p.elmDepth ] = newString(buf, nameStart - bufAbsoluteStart, colonPos - nameStart)
			prefix = p.elPrefix[p.elmDepth]

			// XXX FIX ME FIX ME FIX ME
			// p.elName[ p.elmDepth ] = newString(buf, colonPos + 1 - bufAbsoluteStart, pos - 2 - (colonPos - bufAbsoluteStart))

			name = p.elName[p.elmDepth]
		} else {
			prefix = ""
			p.elPrefix[p.elmDepth] = ""
			// XXX FIX ME
			// p.elName[ p.elmDepth ] = newString(buf, nameStart - bufAbsoluteStart, elLen)
			name = p.elName[p.elmDepth]
		}
	} else {
		// XXX FIX ME FIX ME
		// p.elName[ p.elmDepth ] = newString(buf, nameStart - bufAbsoluteStart, elLen)
		name = p.elName[p.elmDepth]
	}

	for err == nil {
		for u.IsSpace(ch) && err == nil {
			ch, err = p.NextCh()
		} // skip additional white spaces

		if err != nil || ch == '>' {
			break
		} else if ch == '/' {
			// WORKING HERE XXX
			if p.isEmptyElement {
				err = p.NewXmlPullError("repeated / in tag declaration")
			}
			p.isEmptyElement = true
			ch, err = p.NextCh()
			if ch != '>' {
				err = p.NewXmlPullError(
					"expected > to end empty tag not " + printableChar(ch))
			}
			break // XXX inside if ?
		} else if isNameStartChar(ch) {
			// ch = parseAttribute()		// XXX NOT YET !!
			ch, err = p.NextCh()
			continue
		} else {
			err = p.NewXmlPullError(
				"start tag unexpected character " + printableChar(ch))
		}
		//ch, err = p.NextCh(); // skip space
	}

	// now when namespaces were declared we can resolve them
	if p.processNamespaces {
		// uri := getNamespace(prefix)		// XXX NOT YET
		uri := ""
		if uri == "" {
			if prefix == "" { // no prefix and no uri => use default namespace
				uri = NO_NAMESPACE
			} else {
				err = p.NewXmlPullError(
					"could not determine namespace bound to element prefix " + prefix)
			}
		}
		p.elUri[p.elmDepth] = uri

		//String uri = getNamespace(prefix)
		//if uri == nil && prefix == nil) { // no prefix and no uri => use default namespace
		//  uri = ""
		//}
		// resolve attribute namespaces
		for i := 0; i < attributeCount; i++ {
			attrPrefix := p.attributePrefix[i]
			if attrPrefix != "" {
				attrUri := "" // getNamespace(attrPrefix)	// XXX NOT YET
				if attrUri == "" {
					err = p.NewXmlPullError(
						"could not determine namespace bound to attribute prefix " + attrPrefix)
				}
				p.attributeUri[i] = attrUri
			} else {
				p.attributeUri[i] = NO_NAMESPACE
			}
		}
		//[ WFC: Unique Att Spec ]
		// check attribute uniqueness constraint for attributes that has namespace!!!

		for i := 1; i < attributeCount; i++ {
			for j := 0; j < i; j++ {
				if (p.attributeUri[j] == p.attributeUri[i]) &&
					(p.allStringsInterned &&
						(p.attributeName[j] == p.attributeName[i]) ||
						(!p.allStringsInterned &&
							(p.attributeNameHash[j] == p.attributeNameHash[i]) &&
							(p.attributeName[j] == p.attributeName[i]))) {

					// prepare data for nice error message?
					attr1 := p.attributeName[j]
					if p.attributeUri[j] != "" {
						attr1 = p.attributeUri[j] + ":" + attr1
					}
					attr2 := p.attributeName[i]
					if p.attributeUri[i] != "" {
						attr2 = p.attributeUri[i] + ":" + attr2
					}
					err = p.NewXmlPullError(
						"duplicated attributes " + attr1 + " and " + attr2)
				}
			}
		}
	} else { // ! p.processNamespaces

		//[ WFC: Unique Att Spec ]
		// check raw attribute uniqueness constraint!!!
		for i := 1; i < attributeCount; i++ {
			for j := 0; j < i; j++ {
				if p.allStringsInterned &&
					(p.attributeName[j] == p.attributeName[i]) ||
					(!p.allStringsInterned &&
						(p.attributeNameHash[j] == p.attributeNameHash[i]) &&
						(p.attributeName[j] == p.attributeName[i])) {

					// prepare data for nice error message?
					attr1 := p.attributeName[j]
					attr2 := p.attributeName[i]
					err = p.NewXmlPullError(
						"duplicated attributes " + attr1 + " and " + attr2)
				}
			}
		}
	}

	p.elNamespaceCount[p.elmDepth] = p.namespaceEnd
	// posEnd = pos

	_ = name // DEBUG
	_ = prefix

	curEvent = START_TAG
	p.curEvent = curEvent
	return
}
