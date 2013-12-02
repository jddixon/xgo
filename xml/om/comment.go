package om

import (
	"strings"
)

// A comment node, containing text which will not be rendered and may
// be discarded.  Comments are contained within XML comment delimiters
// <pre>
//  &lt;!-- like this --&gt;
// </pre>
//
type Comment struct {
	text string
	Node
}

func NewComment(text string) (this *Comment) {
	node := NewNode()
	text = strings.TrimSpace(text) // XXX Note
	this = &Comment{
		text: text,
		Node: *node,
	}
	return
}

// PROPERTIES ///////////////////////////////////////////////////

// return the text contained within the comment delimiters
func (co *Comment) GetText() string {
	return co.text
}

// NODE METHODS /////////////////////////////////////////////////

// Return true; this node is an Comment.
func (a *Comment) IsComment() bool {
	return true
}

// Return the comment enclosed within XML comment delimiters.
//
func (co *Comment) ToXml() string {

	return "<!-- " + co.text + " -->\n"
}
