package template

import (
	"bufio"
	gc "github.com/jddixon/xgo/context"
	gl "github.com/jddixon/xgo/lex"
	"io"
)

type Template struct {
	ctx    *gc.Context
	rd     io.Reader
	lx     *gl.LexInput
	writer io.Writer // DEBUG ONLY
	wr     *bufio.Writer
}

func NewTemplate(reader io.Reader, writer io.Writer, ctx *gc.Context) (
	t *Template, err error) {

	var lx *gl.LexInput

	if ctx == nil {
		err = NilContext
	} else if reader == nil {
		err = NilReader
	} else if writer == nil {
		err = NilWriter
	} else {
		lx, err = gl.NewNewLexInput(reader) // the doubled New means utf-8
	}
	if err == nil {
		wr := bufio.NewWriter(writer)
		t = &Template{
			ctx:    ctx,
			rd:     reader,
			wr:     wr,
			writer: writer, // DEBUG ONLY
			lx:     lx,
		}
	}
	return
}

func (t *Template) Apply() (s string, err error) {

	var r rune
	for r, err = t.lx.NextCh(); err == nil; r, err = t.lx.NextCh() {
		haveDollar := r == '$'
		if haveDollar {
			r, err = t.lx.NextCh()
			if err != nil {
				break
			} else if r == '{' {
				var sym string
				var value interface{}
				sym, err = t.getSymbol()
				if err != nil {
					break
				}
				value, err = t.ctx.Lookup(sym)
				if err == nil {
					break
				}
				_, err = t.wr.WriteString(value.(string))
			} else {
				_, err = t.wr.WriteRune('$')
				if err == nil {
					_, err = t.wr.WriteRune(r)
				}
			}
		} else {
			_, err = t.wr.WriteRune(r)
		}
	}
	return
}

// Enter having just encountered a ${ sequence with err == nil.
// Consume everything up to the closing right brace, returning the
// symbol as a string.
//
func (t *Template) getSymbol() (s string, err error) {

	var (
		r     rune
		runes []rune
	)
	for r, err = t.lx.NextCh(); err == nil; r, err = t.lx.NextCh() {
		if r == '}' {
			break
		}
		runes = append(runes, r)
	}
	if err == nil {
		s = string(runes)
	}
	return
}
