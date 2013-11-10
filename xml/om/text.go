package om

import (
	"fmt"
)

// An XML Node containing text in String format.
type Text struct {
	text string // CDATA subclass has access
	Node
}

// Create the node, initializing its value.
//
// XXX the text should be XML-escaped.
func NewText(text string) (this *Text) {
	node := NewNode()
	this = &Text{
		text: text,
		Node: *node,
	}
	return
}

func (t *Text) GetText() string {
	return t.text
}

// UNTESTED
func (t *Text) SetText(s string) {
	// XXX SHOULD ESCAPE XML
	t.text = s
}

// NODE METHODS /////////////////////////////////////////////////
// CDATASection overrides.
func (t *Text) IsCdata() bool {
	return false // default
}
func (t *Text) IsText() bool {
	return true
}
func (t *Text) String() string {
	return fmt.Sprintf("[Text:'%s']", t.text)
}

func (t *Text) toXml() string {
	return t.text
}
