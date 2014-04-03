package xmlpull

import (
	"fmt"
	u "unicode"
)

// UTILITIES --------------------------------------------------------

func (p *Parser) checkAttrIndex(index uint) (err error) {
	if index < 0 || index >= p.attributeCount {
		msg := fmt.Sprintf("attribute index must be in 0..%d and not %d",
			p.attributeCount-1, index)
		err = p.NewXmlPullError(msg)
	}
	return
}

// Return a copy of the rune slice.

func MakeCopyRunes(src []rune) (dest []rune, err error) {
	if src == nil || len(src) == 0 {
		err = EmptyRuneSlice
	} else {
		dest = make([]rune, len(src))
		copy(dest, src)
	}
	return
}

func (p *Parser) mustBeStartTag() (err error) {
	if p.curEvent != START_TAG {
		err = p.NewXmlPullError("only START_TAG can have attributes")
	}
	return
} // FOO

// PROPERTIES -------------------------------------------------------
// Sort by the property name, what follows 'get' or 'set'
// ------------------------------------------------------------------

func (p *Parser) getAttributeCount() (n uint) {
	if p.curEvent == START_TAG {
		n = p.attributeCount
	}
	return
}

func (p *Parser) getAttributeName(index uint) (ns string, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		err = p.checkAttrIndex(index)
	}
	if err == nil {
		ns = p.attributeName[index]
	}
	return
}
func (p *Parser) getAttributeNamespace(index uint) (ns string, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		if !p.processNamespaces {
			ns = NO_NAMESPACE
			return // XXX
		} else {
			err = p.checkAttrIndex(index)
		}
	}
	if err == nil {
		ns = p.attributeUri[index]
	}
	return
}
func (p *Parser) getAttributePrefix(index uint) (ns string, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		if p.processNamespaces {
			err = p.checkAttrIndex(index)
		}
		if err == nil {
			ns = p.attributePrefix[index]
		}
	}
	return
}

// XXX NOTICE THAT THE ATTR TYPE RETURNED IS A STRING
func (p *Parser) getAttributeType(index uint) (t string, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		err = p.checkAttrIndex(index)
		if err == nil {
			t = "CDATA"
		}
	}
	return
}

func (p *Parser) getAttributeValue(index uint) (value string, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		err = p.checkAttrIndex(index)
		if err == nil {
			value = p.attributeValue[index]
		}
	}
	return
}

func (p *Parser) getAttributeValueNS(namespace, name string) (
	value string, err error) {

	err = p.mustBeStartTag()
	if err == nil && name == "" {
		err = p.NewXmlPullError("attribute name can not be nil")
	}
	if err == nil {
		if p.processNamespaces {
			for i := uint(0); i < p.attributeCount; i++ {
				if (namespace == p.attributeUri[i] ||
					namespace == p.attributeUri[i]) &&
					name == p.attributeName[i] {
					value = p.attributeValue[i]
				}
			}
		} else {
			if namespace != "" {
				err = p.NewXmlPullError(
					"namespaces processing disabled, attr namespace must be nil")
			} else {
				for i := uint(0); i < p.attributeCount; i++ {
					if name == p.attributeName[i] {
						value = p.attributeValue[i]
					}
				}
			}
		}
	}
	return
}

func (p *Parser) getColumnNumber() int {
	return p.colNo
}

func (p *Parser) getDepth() uint {
	return p.elmDepth
}

func (p *Parser) getEventType() PullEvent {
	return p.curEvent
}

// XXX  Unknown properties are always returned as false
func (p *Parser) getFeature(name string) (found bool, err error) {
	if name == "" {
		err = p.NewXmlPullError("feature name may not be empty")
	} else {
		if FEATURE_PROCESS_NAMESPACES == name {
			found = p.processNamespaces
		} else if FEATURE_NAMES_INTERNED == name {
			found = false
		} else if FEATURE_PROCESS_DOCDECL == name {
			found = false
		} else if FEATURE_XML_ROUNDTRIP == name {
			found = p.roundtripSupported
		}
	}
	return
}

func (p *Parser) setFeature(name string, whether bool) (err error) {

	if name == "" {
		err = p.NewXmlPullError("feature name may not be empty")
	} else {
		if FEATURE_PROCESS_NAMESPACES == name {
			if p.curEvent != START_DOCUMENT {
				err = p.NewXmlPullError(
					"namespace processing feature can only be changed before parsing")
			} else {
				p.processNamespaces = whether
			}
		} else if FEATURE_NAMES_INTERNED == name {
			if whether {
				err = p.NewXmlPullError(
					"interning names in this implementation is not supported")
			}
		} else if FEATURE_PROCESS_DOCDECL == name {
			if whether {
				err = p.NewXmlPullError(
					"processing DOCDECL is not supported")
			}
		} else if FEATURE_XML_ROUNDTRIP == name {
			p.roundtripSupported = whether
		} else {
			msg := fmt.Sprintf("unsupported feature '%s'", name)

			err = p.NewXmlPullError(msg)
		}
	}
	return
}

func (p *Parser) getInputEncoding() string {
	return p.inputEncoding
}

func (p *Parser) getLineNumber() int {
	return p.lineNo
}

// Return a copy of the tag name (in which case the argument is nil
// or empty, or copy the parameter as an entity ref name.
//
// XXX THIS MAKES NO SENSE !

func (p *Parser) getName(candidate string) (s string) {

	if p.curEvent == START_TAG {
		s = p.elName[p.elmDepth]
	} else if p.curEvent == END_TAG {
		s = p.elName[p.elmDepth]
	} else if p.curEvent == ENTITY_REF {
		if p.entityRefName == "" {
			p.entityRefName = candidate
			// XXX ???
			s = candidate
		}
	}
	return
}

// Could/should return error as well
func (p *Parser) getNamespace() (s string) {

	if p.curEvent == START_TAG || p.curEvent == END_TAG {
		if p.processNamespaces {
			s = p.elUri[p.elmDepth]
		}
	}
	return
}

func (p *Parser) getNamespaceCount(elmDepth uint) (n uint, err error) {
	if p.processNamespaces && elmDepth != 0 {
		if elmDepth < 0 || elmDepth > p.elmDepth {
			msg := fmt.Sprintf("elmDepth must be in range 0..%d, but is %d\n",
				p.elmDepth, elmDepth)
			err = p.NewXmlPullError(msg)
		} else {
			n = p.elNamespaceCount[elmDepth]
		}
	}
	return
}

func (p *Parser) getNamespacePrefix(pos uint) (nsP string, err error) {

	if pos < p.namespaceEnd {
		nsP = p.namespacePrefix[pos]
	} else {
		msg := fmt.Sprintf("namespace index %d higher than max", pos)
		err = p.NewXmlPullError(msg)
	}
	return
}

func (p *Parser) getNamespaceFromPrefix(prefix string) (ns string, err error) {

	if prefix != "" {
		for i := p.namespaceEnd - 1; i >= 0; i-- {
			if prefix == p.namespacePrefix[i] {
				ns = p.namespaceUri[i]
			}
		}
		if "xml" == prefix {
			ns = XML_URI
		} else if "xmlns" == prefix {
			ns = XMLNS_URI
		}
	} else {
		for i := p.namespaceEnd - 1; i >= 0; i-- {
			if p.namespacePrefix[i] == "" {
				ns = p.namespaceUri[i]
			}
		}
	}
	return
} // FOO

func (p *Parser) getNamespaceUri(pos uint) (uri string, err error) {
	if pos < p.namespaceEnd {
		uri = p.namespaceUri[pos]
	} else {
		msg := fmt.Sprintf("namespace index %d higher than max", pos)
		err = p.NewXmlPullError(msg)
	}
	return
}

// Return string describing current position of parser:
//   'STATE @line:column'.
//
func (p *Parser) getPositionDescription() (s string) {

	s = fmt.Sprintf("%s @%d:%d",
		PULL_EVENT_NAMES[p.curEvent], p.getLineNumber(), p.getColumnNumber())
	return
}

func (p *Parser) getPrefix() (s string) {

	if p.curEvent == START_TAG {
		s = p.elPrefix[p.elmDepth]
	} else if p.curEvent == END_TAG {
		s = p.elPrefix[p.elmDepth]
	}
	return
}

func (p *Parser) getProperty(name string) (object interface{}, err error) {
	if name == "" {
		err = p.NewXmlPullError("property name must not be empty")
	} else {
		if PROPERTY_XMLDECL_VERSION == name {
			object = p.xmlDeclVersion
		} else if PROPERTY_XMLDECL_STANDALONE == name {
			object = p.xmlDeclStandalone
		} else if PROPERTY_XMLDECL_CONTENT == name {
			object = p.xmlDeclContent
		} else if PROPERTY_LOCATION == name {
			object = p.location
		}
	}
	return
}
func (p *Parser) setProperty(name string, value interface{}) (err error) {
	if PROPERTY_LOCATION == name {
		p.location = value.(string)
	} else {
		msg := fmt.Sprintf("unsupported property: '%s'", name)
		err = p.NewXmlPullError(msg)
	}
	return
}

// XXX NEED TO CHECK ACTUAL USE TO DETERMINE WHETHER THIS MAKES SENSE
//
func (p *Parser) getText() (runes []rune, err error) {
	if p.curEvent != START_DOCUMENT && p.curEvent != END_DOCUMENT {
		if p.curEvent == ENTITY_REF {
			// XXX Why isn't this p.entityRef ???
			runes, err = MakeCopyRunes(p.text)
		} else {
			runes, err = MakeCopyRunes(p.text)
		}
	}
	return
}

func (p *Parser) getTextCharacters(holderForStartAndLength []uint) (
	runes []rune, err error) {

	fmt.Println("getTextCharacters() should never be called")

	// XXX STUB XX
	return
}

// XXX MEANINGLESS
func (p *Parser) isAttributeDefault(index uint) (found bool, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		err = p.checkAttrIndex(index)
		if err == nil {
			found = false
		}
	}
	return
}

func (p *Parser) isEmptyElementTag() (found bool, err error) {

	if p.curEvent != START_TAG {
		err = p.NewXmlPullError(
			"parser must be on START_TAG to check for empty element")
	} else {
		found = p.isEmptyElement
	}
	return
}

// XXX NEED TO CHECK ACTUAL USE TO DETERMINE WHETHER THIS MAKES SENSE
//
func (p *Parser) isWhitespace() (whether bool, err error) {

	if p.curEvent == TEXT || p.curEvent == CDSECT {
		whether = true
		for i := 0; i < len(p.text); i++ {
			if !u.IsSpace(p.text[i]) {
				whether = false
				break
			}
		}
	} else if p.curEvent == IGNORABLE_WHITESPACE {
		whether = true
	} else {
		err = p.NewXmlPullError("no content to check for white spaces")
	}
	return
}

// OTHER PROPERTY-RELATED METHODS ///////////////////////////////////

func (p *Parser) defineEntityReplacementText(
	entityName, replacementText string) (err error) {

	// p.ensureEntityCapacity()

	// make sure that if interning works we take advantage of it

	runes := []rune(entityName)
	p.entityName[p.entityEnd] = entityName
	p.entityReplacement[p.entityEnd] = replacementText
	p.entityNameHash[p.entityEnd] = FastHash(runes)
	p.entityEnd++
	return
}

func (p *Parser) require(_type PullEvent, namespace, name string) (err error) {

	if !p.processNamespaces && namespace != "" {

		err = p.NewXmlPullError(
			"processing namespaces must be enabled to have possible namespaces declared on elements")

	} else if (_type != p.getEventType()) ||
		(namespace != "" && namespace != p.getNamespace()) ||
		(name != "" && (name != p.getName(""))) {

		expectedEvent := PULL_EVENT_NAMES[_type]
		var expectedName, expectedNS string
		if name != "" {
			expectedName = fmt.Sprintf(" with name '%s'", name)
		}
		if namespace != "" {
			expectedNS = fmt.Sprintf(" and namespace '%s'", namespace)
		}
		actualEvent := PULL_EVENT_NAMES[p.getEventType()]
		var actualName, actualNS string
		aName := p.getName("")
		if aName != "" {
			actualName = ", name was '%s'"
		}
		aNS := p.getNamespace()
		if aNS != "" {
			actualNS = ", and namespace was '%s'"
		}

		msg := fmt.Sprintf("expected %s %s %s but actual event was %s %s %s",
			expectedEvent,
			expectedName,
			expectedNS,
			actualEvent, actualName, actualNS)
		err = p.NewXmlPullError(msg)
	}
	return
}

// WORKING HERE; CLEAN UP, SORT, MERGE //////////////////////////////
