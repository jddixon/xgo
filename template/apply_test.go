package template

// xgo/template/apply_test.go

import (
	"fmt"
	gc "github.com/jddixon/xgo_go/context"
	gu "github.com/jddixon/xgo_go/util"
	. "gopkg.in/check.v1"
	"io"
	"strings"
)

var _ = fmt.Print

func (s *XLSuite) TestTemplate(c *C) {

	ctxSs := []string{
		"x\tabc",
		"X\tdef",
	}
	ctxS := strings.Join(ctxSs, "\n") // no terminating newline
	ctx, err := gc.ParseContext(ctxS)
	c.Assert(err, IsNil)
	c.Assert(ctx, NotNil)

	// actually testing context
	c.Assert(ctx.Size(), Equals, len(ctxSs))
	xVal, err := ctx.Lookup("x")
	c.Assert(err, IsNil)
	c.Assert(xVal, Equals, "abc")
	bigXVal, err := ctx.Lookup("X")
	c.Assert(err, IsNil)
	c.Assert(bigXVal, Equals, "def")

	inputSS := []string{
		"this ${x} is ${X} on Sundays",
	}
	expectedSS := []string{
		"this abc is def on Sundays",
	}
	input := strings.Join(inputSS, "\n")
	expected := strings.Join(expectedSS, "\n")

	var rd1 io.Reader = strings.NewReader(input)

	byteBuffer, err := gu.NewByteBuffer(1024)
	c.Assert(err, IsNil)
	c.Assert(byteBuffer, NotNil)
	c.Assert(byteBuffer.Cap(), Equals, 1024)
	c.Assert(byteBuffer.Len(), Equals, 0)

	t, err := NewTemplate(rd1, byteBuffer, ctx)
	c.Assert(err, IsNil)
	c.Assert(t, NotNil)

	err = t.Apply()
	c.Assert(err, Equals, io.EOF)

	output := byteBuffer.String()
	c.Assert(output, Equals, expected)
}
