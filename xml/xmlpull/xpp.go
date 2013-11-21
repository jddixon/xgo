package xmlpull

import (
	"io"
)

// Any fields added here also should be added to reset()
//
type XmlPullParser struct {
	lineNo int // line number
	colNo  int // column number

	// buffer management
	reader   *io.Reader
	encoding string
	buf      []byte
	bufStart int
	bufEnd   int

	// parser state
	curEvent PullToken // defined in const.go

	// element stack
	elmDepth int

	// attribute stack

	// namespace stack

}

func (xpp *XmlPullParser) reset() {
	xpp.lineNo = 1
	xpp.colNo = 0
	xpp.reader = nil
	xpp.encoding = "utf-8"
	xpp.elmDepth = 0
	xpp.bufStart = 0
	xpp.bufEnd = 0

	// XXX STUB XXX

}
func (xpp *XmlPullParser) SetFeature(name string, state bool) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetFeature(name string) (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) SetProperty(name string, value interface{}) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetProperty(name string) (prop interface{}) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) SetInput(in *io.Reader) (err error) {
	if in == nil {
		err = NilReader
	} else {
		xpp.reset()
		xpp.reader = in
	}
	return
}
func (xpp *XmlPullParser) DefineEntityReplacementText(entityName, replacementText string) (err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetNamespaceCount(depth int) (ret int, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetNamespacePrefix(pos int) (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetNamespaceUri(pos int) (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetNamespaceForPrefix(prefix string) (ns string, err error) {
	// XXX STUB XXX
	return
}

//
//
func (xpp *XmlPullParser) GetDepth() int {
	return xpp.elmDepth
}
func (xpp *XmlPullParser) GetPositionDescription() (desc string) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetLineNumber() int {
	return xpp.lineNo
}
func (xpp *XmlPullParser) GetColumnNumber() int {
	return xpp.colNo
}
func (xpp *XmlPullParser) IsWhitespace() (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetText() (t string) {
	// XXX STUB XXX
	return
}

// Should probably return a string.
//
func (xpp *XmlPullParser) GetTextCharacters(holderForStartAndLength []int) (
	chars []byte) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetNamespace() (ns string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetName() (name string) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetPrefix() (prefix string) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) IsEmptyElementTag() (whether bool, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetAttributeCount() (count int) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetAttributeNamespace(index int) (
	ns string, err error) {

	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetAttributeName(index int) (name string) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetAttributePrefix(index int) (p string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetAttributeValue(index int) (val string, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) GetAttributeValueNS(namespace, name string) (
	ns string, err error) {

	// XXX STUB XXX
	return
}

// Return the type of the current parser event.  The spec requires the
// error return.
//
func (xpp *XmlPullParser) GetEventType() (PullToken, error) {
	return xpp.curEvent, nil
}

func (xpp *XmlPullParser) Next() (ret int, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) NextToken() (ret int, err error) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) Require(type_ int, namespace, name string) {
	// XXX STUB XXX
	return
}
func (xpp *XmlPullParser) ReadText() (t string) {
	// XXX STUB XXX
	return
}
