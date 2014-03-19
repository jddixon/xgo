package xmlpull

import (
	"fmt"
	u "unicode"
)

var _ = fmt.Print

// [40] STag ::=  '<' Name (S Attribute)* S? '>'
// [44] EmptyElemTag ::= '<' Name (S Attribute)* S? '/>'

func (p *Parser) parseStartTag() (curEvent PullEvent, err error) {

	// The first character of Name is the next char;  we have seen the
	// opening <
	ch, err := p.NextCh()
	if err != nil {
		return
	}

	var (
		name, prefix []rune
		// elLen        int // XXX DROP ME
	)
	// name = append(name, ch)

	p.elmDepth++
	p.isEmptyElement = false
	attrCount := 0 // so far

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
		// ensureElementsCapacity()			// XXX MPLEMENT ???
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
				p.elPrefix[p.elmDepth] = make([]rune, 0)
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
				// we think we have an attribute
				p.PushBack(ch)
				ch, err = p.parseAttribute()
				// XXX HANDLE ANY ERROR
				// XXX WE SHOULD NOT IGNORE ch

				ch, err = p.NextCh()
				// XXX HANDLE ANY ERROR
				// XXX WE SHOULD NOT IGNORE ch
				continue
			} else {
				err = p.NewXmlPullError(
					"start tag unexpected character " + printableChar(ch))
			}
		}

		// If any namespaces were declared we can now resolve them
		if p.processNamespaces {
			// uri := getNamespace(prefix)		// XXX NOT YET
			var uri []rune // XXX KLUDGE
			if len(uri) == 0 {
				if len(prefix) == 0 { // no prefix and no uri => use default namespace
					uri = NO_NAMESPACE
				} else {
					err = p.NewXmlPullError(
						"can' determine namespace bound to element prefix " +
							string(prefix))
				}
			}
			p.elUri[p.elmDepth] = make([]rune, len(uri))
			copy(p.elUri[p.elmDepth], uri)

			//String uri = getNamespace(prefix)
			//if uri == nil && prefix == nil) { // no prefix and no uri => use default namespace
			//  uri = ""
			//}
			// resolve attribute namespaces
			for i := 0; i < attrCount; i++ {
				attrPrefix := p.attributePrefix[i]
				if len(attrPrefix) > 0 {
					// attrUri := getNamespace(attrPrefix)	// XXX NOT YET
					var attrUri []rune
					if len(attrUri) == 0 {
						err = p.NewXmlPullError(
							"can't determine ns bound to attribute prefix " +
								string(attrPrefix))
					}
					p.attributeUri[i] = make([]rune, len(attrUri))
					copy(p.attributeUri[i], attrUri)
				} else {
					p.attributeUri[i] = NO_NAMESPACE
				}
			}
			//[ WFC: Unique Att Spec ]
			// check attr uniqueness constraint for attrs that have namespace

			for i := 1; i < attrCount; i++ {
				for j := 0; j < i; j++ {
					if SameRunes(p.attributeUri[j], p.attributeUri[i]) &&
						(p.attributeNameHash[j] == p.attributeNameHash[i]) &&
						SameRunes(p.attributeName[j], p.attributeName[i]) {

						// a pretty but rather silly error message
						attr1 := string(p.attributeName[j])
						if len(p.attributeUri[j]) > 0 {
							attr1 = string(p.attributeUri[j]) + ":" + attr1
						}
						attr2 := string(p.attributeName[i])
						if len(p.attributeUri[i]) > 0 {
							attr2 = string(p.attributeUri[i]) + ":" + attr2
						}
						err = p.NewXmlPullError(
							"duplicated attributes " + attr1 + " and " + attr2)
					}
				}
			}
		} else { // ! p.processNamespaces

			//[ WFC: Unique Att Spec ]
			// check raw attribute uniqueness constraint!!!
			for i := 1; i < attrCount; i++ {
				for j := 0; j < i; j++ {
					if SameRunes(p.attributeName[j], p.attributeName[i]) ||
						(p.attributeNameHash[j] == p.attributeNameHash[i]) &&
							SameRunes(p.attributeName[j], p.attributeName[i]) {

						// data for error message
						attr1 := string(p.attributeName[i])
						attr2 := string(p.attributeName[j])
						err = p.NewXmlPullError(
							"duplicated attributes " + attr1 + " and " + attr2)
					}
				}
			}
		}

		p.elNamespaceCount[p.elmDepth] = p.namespaceEnd

		curEvent = START_TAG
		p.curEvent = curEvent

	}
	return
}
