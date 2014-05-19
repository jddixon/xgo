package om

import (
	"strings"
)

// Class representing an XML processing instruction.
//
type PI struct {
	target string // target (language)///
	text   string // text of the instruction///
	Node
}

// Create a node by specifying the target (language) and text
// separately.
//
func NewPI(target string, text string) *PI {

	return &PI{
		target: target,
		text:   text,
	}
}

// Create a node from an initialization string, guessing that
// the first space separates the target from the text.
//
// XXX Needs a better name.
//
func PIFromString(comboText string) (
	this *PI, err error) {

	spaceAt := strings.Index(comboText, " ")
	if spaceAt == -1 {
		err = IllFormedPI
	} else {
		this = &PI{
			target: comboText[0:spaceAt],
			text:   comboText[spaceAt+1:],
		}
	}
	return
}

// PROPERTIES ///////////////////////////////////////////////////

// Return a reference to the target of the PI.
//
func (pi *PI) GetTarget() string {
	return pi.target
}

// Return a reference to the text of the PI.
//
func (pi *PI) GetText() string {
	return pi.text
}

// NODE METHODS /////////////////////////////////////////////////

// Return true; this node is an PI.
func (a *PI) IsPI() bool {
	return true
}

// Output properly bracketed PI content with a line separator,
// without indenting.
//
func (pi *PI) ToXml() string {
	return "<?" + pi.target + " " + pi.text + "?>\n"
}
