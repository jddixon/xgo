package xmlpull

import (
	"fmt"
	u "unicode"
)

// PROPERTIES -------------------------------------------------------
// Sort by the property name, what follows 'get' or 'set'
// ------------------------------------------------------------------

func (xpp *Parser) getColumnNumber() int {
	return xpp.colNo
}

func (xpp *Parser) getDepth() uint {
	return xpp.elmDepth
}

// XXX  Unknown properties are always returned as false
func (xpp *Parser) getFeature(name string) (found bool, err error) {
	if name == "" {
		err = xpp.NewXmlPullError("feature name may not be empty")
	} else {
		if FEATURE_PROCESS_NAMESPACES == name {
			found = xpp.processNamespaces
		} else if FEATURE_NAMES_INTERNED == name {
			found = false
		} else if FEATURE_PROCESS_DOCDECL == name {
			found = false
		} else if FEATURE_XML_ROUNDTRIP == name {
			found = xpp.roundtripSupported
		}
	}
	return
}

func (xpp *Parser) setFeature(name string, whether bool) (err error) {

	if name == "" {
		err = xpp.NewXmlPullError("feature name may not be empty")
	} else {
		if FEATURE_PROCESS_NAMESPACES == name {
			if xpp.curEvent != START_DOCUMENT {
				err = xpp.NewXmlPullError(
					"namespace processing feature can only be changed before parsing")
			} else {
				xpp.processNamespaces = whether
			}
		} else if FEATURE_NAMES_INTERNED == name {
			if whether {
				err = xpp.NewXmlPullError(
					"interning names in this implementation is not supported")
			}
		} else if FEATURE_PROCESS_DOCDECL == name {
			if whether {
				err = xpp.NewXmlPullError(
					"processing DOCDECL is not supported")
			}
		} else if FEATURE_XML_ROUNDTRIP == name {
			xpp.roundtripSupported = whether
		} else {
			msg := fmt.Sprintf("unsupported feature '%s'", name)

			err = xpp.NewXmlPullError(msg)
		}
	}
	return
}

func (xpp *Parser) getInputEncoding() string {
	return xpp.inputEncoding
}

func (xpp *Parser) getLineNumber() int {
	return xpp.lineNo
}

// Return a copy of the tag name (in which case the argument is nil
// or empty, or copy the parameter as an entity ref name.
//
// XXX THIS MAKES NO SENSE !

func (xpp *Parser) getName(candidate string) (s string) {

	if xpp.curEvent == START_TAG {
		s = xpp.elName[xpp.elmDepth]
	} else if xpp.curEvent == END_TAG {
		s = xpp.elName[xpp.elmDepth]
	} else if xpp.curEvent == ENTITY_REF {
		if xpp.entityRefName == "" {
			xpp.entityRefName = candidate
			// XXX ???
			s = candidate
		}
	}
	return
}

// Could/should return error as well
func (xpp *Parser) getNamespace() (s string) {

	if xpp.curEvent == START_TAG || xpp.curEvent == END_TAG {
		if xpp.processNamespaces {
			s = xpp.elUri[xpp.elmDepth]
		}
	}
	return
}

func (xpp *Parser) getNamespaceCount(elmDepth uint) (n uint, err error) {
	if xpp.processNamespaces && elmDepth != 0 {
		if elmDepth < 0 || elmDepth > xpp.elmDepth {
			msg := fmt.Sprintf("elmDepth must be in range 0..%d, but is %d\n",
				xpp.elmDepth, elmDepth)
			err = xpp.NewXmlPullError(msg)
		} else {
			n = xpp.elNamespaceCount[elmDepth]
		}
	}
	return
}

func (xpp *Parser) getNamespacePrefix(pos uint) (nsP string, err error) {

	if pos < xpp.namespaceEnd {
		nsP = xpp.namespacePrefix[pos]
	} else {
		msg := fmt.Sprintf("namespace index %d higher than max", pos)
		err = xpp.NewXmlPullError(msg)
	}
	return
}

func (xpp *Parser) getNamespaceFromPrefix(prefix string) (ns string, err error) {

	if prefix != "" {
		for i := xpp.namespaceEnd - 1; i >= 0; i-- {
			if prefix == xpp.namespacePrefix[i] {
				ns = xpp.namespaceUri[i]
			}
		}
		if "xml" == prefix {
			ns = XML_URI
		} else if "xmlns" == prefix {
			ns = XMLNS_URI
		}
	} else {
		for i := xpp.namespaceEnd - 1; i >= 0; i-- {
			if xpp.namespacePrefix[i] == "" {
				ns = xpp.namespaceUri[i]
			}
		}
	}
	return
} // FOO

func (xpp *Parser) getNamespaceUri(pos uint) (uri string, err error) {
	if pos < xpp.namespaceEnd {
		uri = xpp.namespaceUri[pos]
	} else {
		msg := fmt.Sprintf("namespace index %d higher than max", pos)
		err = xpp.NewXmlPullError(msg)
	}
	return
}

// Return string describing current position of parser:
//   'STATE @line:column'.
//
func (xpp *Parser) getPositionDescription() (s string) {

	s = fmt.Sprintf("%s @%d:%d",
		PULL_EVENT_NAMES[xpp.curEvent], xpp.getLineNumber(), xpp.getColumnNumber())
	return
}

func (xpp *Parser) getPrefix() (s string) {

	if xpp.curEvent == START_TAG {
		s = xpp.elPrefix[xpp.elmDepth]
	} else if xpp.curEvent == END_TAG {
		s = xpp.elPrefix[xpp.elmDepth]
	}
	return
}

func (xpp *Parser) getProperty(name string) (object interface{}, err error) {
	if name == "" {
		err = xpp.NewXmlPullError("property name must not be empty")
	} else {
		if PROPERTY_XMLDECL_VERSION == name {
			object = xpp.xmlDeclVersion
		} else if PROPERTY_XMLDECL_STANDALONE == name {
			object = xpp.xmlDeclStandalone
		} else if PROPERTY_XMLDECL_CONTENT == name {
			object = xpp.xmlDeclContent
		} else if PROPERTY_LOCATION == name {
			object = xpp.location
		}
	}
	return
}
func (xpp *Parser) setProperty(name string, value interface{}) (err error) {
	if PROPERTY_LOCATION == name {
		xpp.location = value.(string)
	} else {
		msg := fmt.Sprintf("unsupported property: '%s'", name)
		err = xpp.NewXmlPullError(msg)
	}
	return
}

// XXX NEED TO CHECK ACTUAL USE TO DETERMINE WHETHER THIS MAKES SENSE
//
func (xpp *Parser) getText() (runes []rune, err error) {
	if xpp.curEvent != START_DOCUMENT && xpp.curEvent != END_DOCUMENT {
		if xpp.curEvent == ENTITY_REF {
			// XXX Why isn't this xpp.entityRef ???
			runes, err = MakeCopyRunes(xpp.text)
		} else {
			runes, err = MakeCopyRunes(xpp.text)
		}
	}
	return
}

func (xpp *Parser) getTextCharacters(holderForStartAndLength []uint) (
	runes []rune, err error) {

	fmt.Println("getTextCharacters() should never be called")

	// XXX STUB XX
	return
}

// XXX MEANINGLESS
func (xpp *Parser) isAttributeDefault(index uint) (found bool, err error) {
	err = xpp.mustBeStartTag()
	if err == nil {
		err = xpp.checkAttrIndex(index)
		if err == nil {
			found = false
		}
	}
	return
}

func (xpp *Parser) isEmptyElementTag() (found bool, err error) {

	if xpp.curEvent != START_TAG {
		err = xpp.NewXmlPullError(
			"parser must be on START_TAG to check for empty element")
	} else {
		found = xpp.isEmptyElement
	}
	return
}

// XXX NEED TO CHECK ACTUAL USE TO DETERMINE WHETHER THIS MAKES SENSE
//
func (xpp *Parser) isWhitespace() (whether bool, err error) {

	if xpp.curEvent == TEXT || xpp.curEvent == CDSECT {
		whether = true
		for i := 0; i < len(xpp.text); i++ {
			if !u.IsSpace(xpp.text[i]) {
				whether = false
				break
			}
		}
	} else if xpp.curEvent == IGNORABLE_WHITESPACE {
		whether = true
	} else {
		err = xpp.NewXmlPullError("no content to check for white spaces")
	}
	return
}

// OTHER PROPERTY-RELATED METHODS ///////////////////////////////////

func (xpp *Parser) defineEntityReplacementText(
	entityName, replacementText string) (err error) {

	// xpp.ensureEntityCapacity()

	// make sure that if interning works we take advantage of it

	runes := []rune(entityName)
	xpp.entityName[xpp.entityEnd] = entityName
	xpp.entityReplacement[xpp.entityEnd] = replacementText
	xpp.entityNameHash[xpp.entityEnd] = FastHash(runes)
	xpp.entityEnd++
	return
}
