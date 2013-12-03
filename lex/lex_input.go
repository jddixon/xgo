package lex

// xgo/xml/lex/lex_input.go

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	u "unicode"
)

const (
	MAX_PUSH_BACK = 256 // an arbitrary number
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

// Expect the next character to match the expected char,
// returning an error if there is a mismatch.
func (lx *LexInput) ExpectCh(expected rune) (err error) {

	ch, err := lx.NextCh()
	if err == nil && ch != expected {
		msg := fmt.Sprintf("expected '%c', found '%c'", expected, ch)
		err = errors.New(msg)
		lx.PushBack(ch)
	}
	return
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

// Expect the next characters to exactly match the expected string,
// returning an error if there is a mismatch.  If the parameter is
// an empty string, do nothing.
func (lx *LexInput) ExpectStr(expected string) (err error) {

	for _, r := range expected {
		var ch rune
		ch, err = lx.NextCh()
		if err == nil && ch != r {
			msg := fmt.Sprintf("expected '%c' in '%s', found '%c'",
				r, expected, ch)
			err = errors.New(msg)
		}
		if err != nil {
			break
		}
	}
	return
}

// If the next characters exactly match the parameter, return true.
// If there is a mismatch, return false and push any characters read
// back on the input.
func (lx *LexInput) AcceptStr(acceptable string) (found bool, err error) {

	var runes []rune
	found = true
	for _, r := range acceptable {
		var ch rune
		ch, err = lx.NextCh()
		if err != nil {
			break
		}
		runes = append(runes, ch)
		if ch != r {
			found = false
			break
		}
	}
	// If we didn't get a match, we need to push back any characters
	// read, FIFO-style.
	if err == nil && !found {
		for i := 0; i < len(runes); i++ {
			lx.PushBack(runes[i])
		}
	}
	return
}
