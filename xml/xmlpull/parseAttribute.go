package xmlpull

// parseAttribute.go

import (
	"fmt"
)

var _ = fmt.Print

// [41] Attribute ::= Name Eq AttValue
// [WFC: No External Entity References]
// [WFC: No < in Attribute Values]
//
// XXX WHAT CHAR IS RETURNED?
//
func (p *Parser) parseAttribute() (ch rune, err error) {

	ch, err = p.NextCh()
	// XXX HANDLE ERROR
	if ch == ':' && p.processNamespaces {
		err = p.NewXmlPullError(
			"namespaces processing enabled; colon can't be at attr name start")
		return // DIJKSTRA
	}

	var cName []rune
	cName = append(cName, ch) // first char in attr cName

	// we test to see whether this is actually true
	startsWithXmlns := p.processNamespaces && ch == 'x'
	var colonPos uint // we rely on colon being forbidden as first char
	var xmlnsPos uint

	ch, err = p.NextCh()
	for err == nil && isNameChar(ch) {
		if p.processNamespaces {
			if startsWithXmlns && xmlnsPos < 5 {
				xmlnsPos++
				if xmlnsPos == 1 {
					if ch != 'm' {
						startsWithXmlns = false
					}
				} else if xmlnsPos == 2 {
					if ch != 'l' {
						startsWithXmlns = false
					}
				} else if xmlnsPos == 3 {
					if ch != 'n' {
						startsWithXmlns = false
					}
				} else if xmlnsPos == 4 {
					if ch != 's' {
						startsWithXmlns = false
					}
				} else if xmlnsPos == 5 {
					if ch != ':' {
						err = p.NewXmlPullError(
							"colon must follow 'xmlns' in attr cName " +
								"when namespaces are enabled")
						break
					}
				}
			}
			if ch == ':' {
				if colonPos > 0 {
					err = p.NewXmlPullError(
						"only one colon is allowed in attribute cName when namespaces are enabled")
					break
				}
				colonPos = uint(len(cName))
			}
		}
		cName = append(cName, ch)
		ch, err = p.NextCh()
	}
	if err != nil {
		return
	}

	// XXX Kukemal !
	// ensureAttributesCapacity(p.attributeCount)

	// var prefixRunes []rune
	var prefix string
	// var prefixLen uint
	var name string
	var nameRunes []rune
	var nameLen uint
	cNameLen := uint(len(cName))

	// work on prefixes and namespace URI
	if p.processNamespaces {
		if xmlnsPos < 4 {
			startsWithXmlns = false
		}
		if startsWithXmlns {
			if colonPos > 0 {
				nameLen = cNameLen - colonPos
				if nameLen == 0 {
					err = p.NewXmlPullError(
						"namespace prefix is required after xmlns: when namespaces are enabled")
					return // DIJKSTRA
				}
				name = string(cName[colonPos+1:])
			}
		} else {
			if colonPos > 0 {
				// prefixLen = colonPos - 1
				prefix = string(cName[:colonPos])
				p.attributePrefix[p.attributeCount] = prefix

				// XXX THIS IS STUPID: assert calculated len = real len
				nameLen = cNameLen - colonPos
				name = string(cName[colonPos+1:])
			} else {
				prefix = ""
				p.attributePrefix[p.attributeCount] = ""

				nameLen = cNameLen
				name = string(cName)
				p.attributeName[p.attributeCount] = name
			}

			// FashHash replaces Java String.hashCode
			p.attributeNameHash[p.attributeCount] = FastHash(nameRunes)
		}

	} else {
		nameLen = cNameLen
		nameRunes = cName
		name = string(nameRunes)
		p.attributeName[p.attributeCount] = name

		// FashHash replaces Java String.hashCode
		p.attributeNameHash[p.attributeCount] = FastHash(nameRunes)
	}

	// XXX WORKING HERE XXX

	// [25] Eq ::=  S? '=' S?

	// skip any spaces
	p.SkipS()

	ch, err = p.NextCh()
	if err == nil {
		if ch != '=' {
			err = p.NewXmlPullError("expected = after attribute name")
		}
		p.SkipS()
	}
	if err != nil {
		return
	}

	// [10] AttValue ::=   '"' ([^<&"] | Reference)* '"'
	//                  |  "'" ([^<&'] | Reference)* "'"

	DELIM := ch // consider me a constant, please

	if DELIM != '"' && DELIM != '\'' {
		msg := fmt.Sprintf(
			"attr value must start with quotation or apostrophe not '%c'",
			DELIM)
		err = p.NewXmlPullError(msg)
		return
	}

	// parse until DELIM or < and resolve Reference
	//[67] Reference ::= EntityRef | CharRef

	normalizedCR := false // NOT USED

	var valueRunes []rune
	var value string
	p.entityRefName = p.entityRefName[:0]

	for err == nil {
		ch, err = p.NextCh()
		if err != nil {
			return
		}
		if ch == DELIM {
			break
		}
		if ch == '<' {
			err = p.NewXmlPullError(
				"found < : markup not allowed inside attribute value")
			return
		}
		if ch == '&' {
			// extractEntityRef
			// XXX DESPERATELY NEEDS TESTING
			var resolvedEntity []rune
			resolvedEntity, err = p.parseEntityRef()
			// check if replacement text can be resolved !!!
			if resolvedEntity == nil {
				if p.entityRefName == "" {
					// XXX SHOULD BE A STRING OR IT SHD BE COPIED
					p.entityRefName = string(valueRunes)
				}
				msg := fmt.Sprintf("Could not resolve entity '%s'",
					p.entityRefName)
				err = p.NewXmlPullError(msg)
				return
			}
			valueRunes = append(valueRunes, resolvedEntity...)

		} else if ch == '\t' || ch == '\n' || ch == '\r' {
			// XXX DESPERATELY NEEDS TESTING
			// do attribute value normalization
			// as described in http://www.w3.org/TR/REC-xml#AVNormalize
			// handle EOL normalization ...
			valueRunes = append(valueRunes, ' ')
		} else {
			valueRunes = append(valueRunes, ch)
		}
		normalizedCR = ch == '\r'
	}
	_ = normalizedCR // NOT USED
	value = string(valueRunes)

	if p.processNamespaces && startsWithXmlns {

		// XXX WE ARE NOT CAPTURING THIS CORRECTLY
		ns := string(cName)
		nsLen := cNameLen

		// ensureNamespacesCapacity(namespaceEnd)

		var prefixHash uint32

		if colonPos > 0 {
			if nsLen == 0 {
				err = p.NewXmlPullError(
					"non-default namespace can not be declared to be empty string")
				return
			}
			// declare new namespace
			p.namespacePrefix[p.namespaceEnd] = string(name)
			prefixHash = FastHash(nameRunes)
			p.namespacePrefixHash[p.namespaceEnd] = prefixHash
		} else {
			// declare  new default namespace ...
			p.namespacePrefix[p.namespaceEnd] = ""
			prefixHash = 0 // XXX MAJOR PROBLEM
			p.namespacePrefixHash[p.namespaceEnd] = prefixHash
		}
		p.namespaceUri[p.namespaceEnd] = ns

		// detect any duplicate namespace declarations
		startNs := p.elNamespaceCount[p.elmDepth-1]
		for i := p.namespaceEnd - 1; i >= startNs; i-- {
			if nameLen > 0 &&
				p.namespacePrefixHash[i] == prefixHash &&
				ns == p.namespacePrefix[i] {

				msg := fmt.Sprintf(
					"duplicated namespace declaration for '%s' prefix", ns)
				err = p.NewXmlPullError(msg)
				return
			}
		}

		p.namespaceEnd++

	} else {
		// XXX NEEDS TESTING
		p.attributeValue[p.attributeCount] = value
		p.attributeCount++
	}
	return
}
