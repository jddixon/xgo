package xmlpull

import (
	"fmt"
	"io"
)

// UTILITIES --------------------------------------------------------

func (xpp *Parser) checkAttrIndex(index uint) (err error) {
	if index < 0 || index >= xpp.attributeCount {
		msg := fmt.Sprintf("attribute index must be in 0..%d and not %d",
			xpp.attributeCount-1, index)
		err = xpp.NewXmlPullError(msg)
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

func (xpp *Parser) mustBeStartTag() (err error) {
	if xpp.curEvent != START_TAG {
		err = xpp.NewXmlPullError("only START_TAG can have attributes")
	}
	return
}

/////////////////////////////////////////////////////////////////////
// This code used to live in parser.go.  There are 30 functions here.
/////////////////////////////////////////////////////////////////////

// ==================================================================
// IMPLMENTATION OF XmlPullParserI
// ==================================================================

func (xpp *Parser) SetFeature(name string, state bool) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetFeature(name string) (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) SetProperty(name string, value interface{}) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetProperty(name string) (prop interface{}) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) SetInput(in *io.Reader) (err error) {
	if in == nil {
		err = NilReader
	} else {
		xpp.reset()
		xpp.LexInput.Reset()
	}
	return
}
func (xpp *Parser) DefineEntityReplacementText(entityName, replacementText string) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespaceCount(depth int) (ret int, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespacePrefix(pos int) (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespaceUri(pos int) (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespaceForPrefix(prefix string) (ns string, err error) {
	// XXX STUB XXX
	return
}

func (xpp *Parser) GetDepth() uint {
	return xpp.elmDepth
}
func (xpp *Parser) GetPositionDescription() (desc string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetLineNumber() int {
	return xpp.lineNo
}
func (xpp *Parser) GetColumnNumber() int {
	return xpp.colNo
}
func (xpp *Parser) IsWhitespace() (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetText() (t string) {
	// XXX STUB XXX
	return
}

// Should probably return a string.
//
func (xpp *Parser) GetTextCharacters(holderForStartAndLength []int) (
	chars []byte) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetNamespace() (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetName() (name string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetPrefix() (prefix string) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) IsEmptyElementTag() (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *Parser) GetAttributeCount() (n uint) {
	if xpp.curEvent == START_TAG {
		n = xpp.attributeCount
	}
	return
}

func (xpp *Parser) GetAttributeName(index uint) (ns string, err error) {
	err = xpp.mustBeStartTag()
	if err == nil {
		err = xpp.checkAttrIndex(index)
	}
	if err == nil {
		ns = xpp.attributeName[index]
	}
	return
}
func (xpp *Parser) GetAttributeNamespace(index uint) (ns string, err error) {
	err = xpp.mustBeStartTag()
	if err == nil {
		if !xpp.processNamespaces {
			ns = NO_NAMESPACE
			return // XXX
		} else {
			err = xpp.checkAttrIndex(index)
		}
	}
	if err == nil {
		ns = xpp.attributeUri[index]
	}
	return
}
func (xpp *Parser) GetAttributePrefix(index uint) (ns string, err error) {
	err = xpp.mustBeStartTag()
	if err == nil {
		if xpp.processNamespaces {
			err = xpp.checkAttrIndex(index)
		}
		if err == nil {
			ns = xpp.attributePrefix[index]
		}
	}
	return
}
func (xpp *Parser) GetAttributeType(index uint) (t string, err error) {
	err = xpp.mustBeStartTag()
	if err == nil {
		err = xpp.checkAttrIndex(index)
		if err == nil {
			t = "CDATA"
		}
	}
	return
}

func (xpp *Parser) GetAttributeValue(index uint) (value string, err error) {
	err = xpp.mustBeStartTag()
	if err == nil {
		err = xpp.checkAttrIndex(index)
		if err == nil {
			value = xpp.attributeValue[index]
		}
	}
	return
}
func (xpp *Parser) GetAttributeValueNS(namespace, name string) (
	value string, err error) {

	err = xpp.mustBeStartTag()
	if err == nil && name == "" {
		err = xpp.NewXmlPullError("attribute name can not be nil")
	}
	if err == nil {
		if xpp.processNamespaces {
			for i := uint(0); i < xpp.attributeCount; i++ {
				if (namespace == xpp.attributeUri[i] ||
					namespace == xpp.attributeUri[i]) &&
					name == xpp.attributeName[i] {
					value = xpp.attributeValue[i]
				}
			}
		} else {
			if namespace != "" {
				err = xpp.NewXmlPullError(
					"namespaces processing disabled, attr namespace must be nil")
			} else {
				for i := uint(0); i < xpp.attributeCount; i++ {
					if name == xpp.attributeName[i] {
						value = xpp.attributeValue[i]
					}
				}
			}
		}
	}
	return
}

// Return the type of the current parser event.  XXX The spec requires an
// error return.
//
func (xpp *Parser) GetEventType() PullEvent {
	return xpp.curEvent
}

func (xpp *Parser) Require(_type PullEvent, namespace, name string) (
	err error) {

	if !xpp.processNamespaces && namespace != "" {

		err = xpp.NewXmlPullError(
			"processing namespaces not enabled but namespaces declared on elements")

	} else if (_type != xpp.GetEventType()) ||
		(namespace != "" && namespace != xpp.getNamespace()) ||
		(name != "" && (name != xpp.getName(""))) {

		expectedEvent := PULL_EVENT_NAMES[_type]
		var expectedName, expectedNS string
		if name != "" {
			expectedName = fmt.Sprintf(" with name '%s'", name)
		}
		if namespace != "" {
			expectedNS = fmt.Sprintf(" and namespace '%s'", namespace)
		}
		actualEvent := PULL_EVENT_NAMES[xpp.GetEventType()]
		var actualName, actualNS string
		aName := xpp.getName("")
		if aName != "" {
			actualName = ", name was '%s'"
		}
		aNS := xpp.getNamespace()
		if aNS != "" {
			actualNS = ", and namespace was '%s'"
		}

		msg := fmt.Sprintf("expected %s %s %s but actual event was %s %s %s",
			expectedEvent,
			expectedName,
			expectedNS,
			actualEvent, actualName, actualNS)
		err = xpp.NewXmlPullError(msg)
	}
	return
}

func (xpp *Parser) ReadText() (t string) {
	// XXX STUB XXX
	return
}
