package om

import ()

// An XML Node containing text in String format.
type Text struct {
	text string // CDATA subclass has access
	Node
}

// Create the node, initializing its value.
//
// XXX The text should be XML-escaped.
// XXX Should the text be trimmed?
func NewText(text string) *Text {
	node := NewNode()
	return &Text{
		text: text,
		Node: *node,
	}
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

// Return true; this node is an Text.
func (a *Text) IsText() bool {
	return true
}

func (t *Text) ToXml() string {
	return t.text
}
