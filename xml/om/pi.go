package om

import (
	"strings"
)

// Class representing an XML processing instruction.
//
type ProcessingInstruction struct {
	target string // target (language)///
	text   string // text of the instruction///
	Node
}

// Create a node by specifying the target (language) and text
// separately.
//
func NewProcessingInstruction(target string, text string) (
	*ProcessingInstruction) {

	return &ProcessingInstruction {
		target : target,
		text : text,
	}
}

// Create a node from an initialization string, guessing that
// the first space separates the target from the text.
//
// XXX Needs a better name.
//
func ProcessingInstructionFromString(comboText string) (
	this *ProcessingInstruction, err error) {

	spaceAt := strings.Index(comboText, " ")
	if spaceAt == -1 {
		err = IllFormedPI
	} else {
		this = & ProcessingInstruction {
			target:  comboText[0:spaceAt],
			text: comboText[spaceAt+1:],
		}
	}
	return
}

// PROPERTIES ///////////////////////////////////////////////////

// Return a reference to the target of the PI.
//
func (pi *ProcessingInstruction) GetTarget() string {
	return pi.target
}

// Return a reference to the text of the PI.
//
func (pi *ProcessingInstruction) GetText() string {
	return pi.text
}

// NODE METHODS /////////////////////////////////////////////////

// Output properly bracketed PI content with a line separator,
// without indenting.
//
func (pi *ProcessingInstruction) ToXml() string {
	return "<?" + pi.target + " " + pi.text + "?>\n"
}
