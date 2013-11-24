package lex

// xgo/xml/lex/lex_input.go

import (
	"io"
)

type LexInput struct {
	encoding string // should default to "utf-8"
	reader   *io.Reader

	lineNo int // line number, 1-based
	colNo  int // column number, 1-based but initially zero
	Buf    []byte
	Offset int
	End    int
}

func (lx *LexInput) Reset() {
	lx.encoding = "utf-8"
	lx.lineNo = 1
	lx.colNo = 0
	lx.Buf = nil // questionable
	lx.Offset = 0
	lx.End = 0
}

func (lx *LexInput) Encoding() string {
	return lx.encoding
}

func (lx *LexInput) SetInput(in *io.Reader) {
	lx.Reset()
	lx.reader = in
}
func (lx *LexInput) LineNo() int {
	return lx.lineNo
}
func (lx *LexInput) stepLineNo() int {
	lx.lineNo++
	lx.colNo = 0
	return lx.lineNo
}
func (lx *LexInput) ColNo() int {
	return lx.colNo
}

func (lx *LexInput) SkipS() {
	// XXX Convert to use runes.
	for (lx.Offset < lx.End) && IsS(lx.Buf[lx.Offset]) {
		lx.Offset++
	}
}

func (lx *LexInput) ExpectS() (err error) {
	if !IsS(lx.Buf[lx.Offset]) {
		err = ExpectedSpace
	} else {
		lx.SkipS()
	}
	return
}

// XXX Drop this, use unicode/IsSpace(r rune) bool (and similar functions)
// instead
func IsS(ch byte) bool {
	return ch == 0x20 || ch == 0x0a || ch == 0x0d || ch == 0x09
}
