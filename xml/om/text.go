package om

import ()

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

// Return true; this node is an Text.
func (a *Text) IsText() bool {
	return true
}

func (t *Text) ToString() string {
	return "[Text:'" + t.text + "']"
}

func (t *Text) ToXml() string {
	return t.text
}
