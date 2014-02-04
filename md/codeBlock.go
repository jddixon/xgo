package md

// xgo/md/code_block.go

import (
	"fmt"
)

var _ = fmt.Print

// A CodeBlock is a Block consisting of a number of lines, collected
// here as CodeLines
type CodeBlock struct {
	lines []CodeLine // to force conversion to entities
}

// XXX Might make more sense to copy.
func (p *CodeBlock) Add(line *CodeLine) {
	p.lines = append(p.lines, *line)
	return
}

func (p *CodeBlock) Clear() {
	p.lines = p.lines[:0]
}

func (p *CodeBlock) GetHtml() (runes []rune) {
	runes = append(runes, CODE_OPEN...)
	for i := 0; i < len(p.lines); i++ {
		text := p.lines[i].GetHtml()
		runes = append(runes, text...)
		runes = append(runes, LF)
	}
	runes = append(runes, CODE_CLOSE...)
	return
}

func (p *CodeBlock) Size() uint {
	return uint(len(p.lines))
}