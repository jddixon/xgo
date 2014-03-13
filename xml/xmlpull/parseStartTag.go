package xmlpull

import (
	"fmt"
	u "unicode"
)

var _ = fmt.Print

// [40] STag ::=  '<' Name (S Attribute)* S? '>'
// [44] EmptyElemTag ::= '<' Name (S Attribute)* S? '/>'

func (p *Parser) parseStartTag(ch rune) (curEvent PullEvent, err error) {

	// The first character of Name is the parameter, so we have seen the
	// opening <
	var (
		name, prefix []rune
		elLen        int // XXX DROP ASAP
	)
	name = append(name, ch)

	p.elmDepth++
	p.isEmptyElement = false
	attributeCount := 0 // so far

	var colonFound bool
	if ch == ':' && p.processNamespaces {
		err = p.NewXmlPullError(
			"ns processing enabled: colon cannot start element name")
		// DIJKSTRA
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
					"ns processing enabled: only one colon allowed in name of element")
				break
			}
			colonFound = true
		}
		if ch == ':' {
			prefix = make([]rune, len(name))
			copy(prefix, name)
			name = name[:0]
		} else {
			name = append(name, ch)
		}
	}
	// we have a name and may have a prefix
	if err == nil {
		// ensureElementsCapacity()			// XXX MUST IMPLEMENT
		elLen = len(name) // XXX useless
		p.elRawName[p.elmDepth] = name
		p.elRawNameLine[p.elmDepth] = p.lineNo

		// work on prefixes and namespace URI
		if p.processNamespaces {
			if colonFound {
				// XXX FIX ME FIX ME FIX ME
				// p.elPrefix[ p.elmDepth ] = newString(buf, nameStart - bufAbsoluteStart, colonPos - nameStart)
				prefix = p.elPrefix[p.elmDepth] // XXX WRONG

				// XXX FIX ME FIX ME FIX ME
				// p.elName[ p.elmDepth ] = newString(buf, colonPos + 1 - bufAbsoluteStart, pos - 2 - (colonPos - bufAbsoluteStart))

				name = p.elName[p.elmDepth]
			} else {
				// prefix is empty
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

	} // FOO
	return
}
