package lex

// xgo/xml/lex/lex_input.go

import (
	"bufio"
	"io"
	u "unicode"
)

const (
	MAX_PUSH_BACK = 3 // an arbitrary number
)

type LexInput struct {
	encoding   string // should default to "utf-8"
	rd         *bufio.Reader
	lineNo     int // line number, 1-based
	colNo      int // column number, character-based, initially zero
	offset     int // byte offset into the source
	pushedBack []rune
	End        int
}

// An io.Reader has a Read(p []byte) (int, error) function that reads
// up to len(p) bytes into p, returning the number of bytes read.  If
// Read() encounters an error or EOF, it returns the number of bytes
// successfully read and an error code, so the rule is to process the
// bytes and then look at the error.  If Reader() gets an EOF, it will
// return the number of bytes read and either nil or EOF.  If it
// returns nil, it will return 0, EOF on the next read.
//
func NewLexInput(reader io.Reader, encoding string) (
	lx *LexInput, err error) {

	if reader == nil {
		err = NilReader
	} else {
		lx = new(LexInput)
		lx.Reset()
		lx.rd = bufio.NewReader(reader)
		if encoding != "" {
			// XXX Need some validation
			lx.encoding = encoding
		}
	}
	return
}

func NewNewLexInput(reader io.Reader) (*LexInput, error) {
	return NewLexInput(reader, "utf-8")
}

// Return the byte offset into the input.
//
func (lx *LexInput) GetOffset() int {
	return lx.offset
}

// Read the next character (rune) on the input, adjusting the offset
// accordingly.  The column number is incremented, but using software is
// responsible for dealing with line breaks.
//
func (lx *LexInput) NextCh() (r rune, err error) {
	if len(lx.pushedBack) > 0 {
		r = lx.pushedBack[0]
		lx.pushedBack = lx.pushedBack[1:]
	} else {
		var delta int
		r, delta, err = lx.rd.ReadRune()
		if delta > 0 {
			lx.offset += delta // byte offset
			lx.colNo++         // character offset
		}
	}
	return
}

// Initialize the data structure.
//
func (lx *LexInput) Reset() {
	var pb []rune
	lx.encoding = "utf-8"
	lx.lineNo = 1
	lx.colNo = 0
	lx.offset = 0
	lx.End = 0
	lx.pushedBack = pb
}

func (lx *LexInput) Encoding() string {
	return lx.encoding
}

// XXX This function introduces complexity which may not be dealt with.
//
func (lx *LexInput) SetInput(in io.Reader) {
	lx.Reset()
	lx.rd = bufio.NewReader(in)
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

// Push back the character r onto the input.  If this exceeds the max,
// the character is silently discarded.
//
// XXX This treatment may not be acceptable.
//
func (lx *LexInput) PushBack(r rune) {
	if len(lx.pushedBack) < MAX_PUSH_BACK {
		lx.pushedBack = append(lx.pushedBack, r)
		lx.offset--
		// XXX Column numbers are not handled!
		// XXX Need to be able to push back newlines - which affect
		//     colNo, lineNo, and offset.
	}
}

// Skip spaces on the input.  If an error is encountered on reading the
// next character, it is pushed back and the error is silently discarded.
//
// XXX This treatment may not be acceptable.
//
func (lx *LexInput) SkipS() {
	for {
		r, err := lx.NextCh()
		if err != nil {
			lx.PushBack(r)
		} else {
			if !u.IsSpace(r) {
				lx.pushedBack = append(lx.pushedBack, r)
				lx.offset--
				break
			}
		}
	}
}

// Expect the next char to be a space; if it isn't, return an error.
// Otherwise, skip any spaces found
//
func (lx *LexInput) ExpectS() (err error) {
	r, err := lx.NextCh()
	if err != nil {
		lx.PushBack(r)
	} else {
		if !u.IsSpace(r) {
			err = ExpectedSpace
		} else {
			lx.SkipS()
		}
	}
	return
}

// XXX Drop this, use unicode/IsSpace(r rune) bool (and similar functions)
// instead
func IsS(ch byte) bool {
	return ch == 0x20 || ch == 0x0a || ch == 0x0d || ch == 0x09
}
