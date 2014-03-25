package xmlpull

import (
	"fmt"
)

// UTILITIES --------------------------------------------------------

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

// PROPERTIES -------------------------------------------------------

func (p *Parser) getAttributeCount() int {
	if p.curEvent != START_TAG {
		return -1
	}
	return p.attributeCount
}

// Return a copy of the tag name (in which case the argument is nil
// or empty, or copy the parameter as an entity ref name.

func (p *Parser) getName(candidate []rune) (runes []rune, err error) {

	if p.curEvent == START_TAG {
		runes, err = MakeCopyRunes(p.elName[p.elmDepth])
	} else if p.curEvent == END_TAG {
		runes, err = MakeCopyRunes(p.elName[p.elmDepth])
	} else if p.curEvent == ENTITY_REF {
		if p.entityRefName == nil {
			runes, err = MakeCopyRunes(candidate)
			if err == nil {
				p.entityRefName, err = MakeCopyRunes(runes)
			}
		}
	}
	if err != nil {
		runes = nil
	}
	return
}

func (p *Parser) getNamespace() (s string, err error) {

	if p.curEvent == START_TAG || p.curEvent == END_TAG {
		if p.processNamespaces {
			s = p.elUri[p.elmDepth]
		}
	}
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

func (p *Parser) isEmptyElementTag() (found bool, err error) {

	if p.curEvent != START_TAG {
		err = p.NewXmlPullError(
			"parser must be on START_TAG to check for empty element")
	} else {
		found = p.isEmptyElement
	}
	return
}

// WORKING HERE; CLEAN UP, SORT, MERGE

// UTILITY
func (p *Parser) checkAttrIndex(index int) (err error) {
	if index < 0 || index >= p.attributeCount {
		msg := fmt.Sprintf("attribute index must be in 0..%d and not %d",
			p.attributeCount-1, index)
		err = p.NewXmlPullError(msg)
	}
	return
}
func (p *Parser) mustBeStartTag() (err error) {
	if p.curEvent != START_TAG {
		err = p.NewXmlPullError("only START_TAG can have attributes")
	}
	return
}

// PROPERTIES

func (p *Parser) getAttributeNamespace(index int) (ns string, err error) {
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
func (p *Parser) getAttributeName(index int) (ns string, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		err = p.checkAttrIndex(index)
	}
	if err == nil {
		ns = p.attributeName[index]
	}
	return
}
func (p *Parser) getAttributePrefix(index int) (ns string, err error) {
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
func (p *Parser) getAttributeType(index int) (t string, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		err = p.checkAttrIndex(index)
		if err == nil {
			t = "CDATA"
		}
	}
	return
}

// XXX MEANINGLESS
func (p *Parser) isAttributeDefault(index int) (found bool, err error) {
	err = p.mustBeStartTag()
	if err == nil {
		err = p.checkAttrIndex(index)
		if err == nil {
			found = false
		}
	}
	return
}

func (p *Parser) getAttributeValue(index int) (value string, err error) {
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
			for i := 0; i < p.attributeCount; i++ {
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
				for i := 0; i < p.attributeCount; i++ {
					if name == p.attributeName[i] {
						value = p.attributeValue[i]
					}
				}
			}
		}
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
}
