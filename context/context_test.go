package context

// xgo/context/context_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

func (s *XLSuite) TestEmpty(c *C) {
	ctx := NewNewContext()
	c.Assert(ctx, NotNil)
	c.Assert(ctx.Size(), Equals, 0)
	c.Assert(ctx.GetParent(), IsNil)
	value, err := ctx.Lookup("foo")
	c.Assert(err, IsNil)
	c.Assert(value, IsNil)
}

func (s *XLSuite) TestAddingNulls(c *C) {
	ctx := NewNewContext()
	err := ctx.Bind("", "bar")
	c.Assert(err, NotNil)
	err = ctx.Bind("foo", nil)
	c.Assert(err, NotNil)
}
func (s *XLSuite) TestSimpleBindings(c *C) {
	ctx := NewNewContext()
	err := ctx.Bind("foo", "that was foo")
	c.Assert(err, IsNil)
	err = ctx.Bind("bar", "that was bar")
	c.Assert(err, IsNil)
	c.Assert(ctx.Size(), Equals, 2)
	value, err := ctx.Lookup("foo")
	c.Assert(value, Equals, "that was foo")
	value, err = ctx.Lookup("bar")
	c.Assert(value, Equals, "that was bar")
}
func (s *XLSuite) TestNestedContexts(c *C) {
	ctx := NewNewContext()
	ctx1 := NewContext(ctx)
	ctx2 := NewContext(ctx1)
	c.Assert(ctx1.GetParent(), Equals, ctx)
	c.Assert(ctx2.GetParent(), Equals, ctx1)
	err := ctx.Bind("foo", "bar0")
	c.Assert(err, IsNil)
	err = ctx1.Bind("foo", "bar1")
	c.Assert(err, IsNil)
	err = ctx2.Bind("foo", "bar2")
	c.Assert(err, IsNil)

	value, err := ctx2.Lookup("foo")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "bar2")

	err = ctx2.Unbind("foo")
	c.Assert(err, IsNil)
	value, err = ctx2.Lookup("foo")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "bar1")

	err = ctx1.Unbind("foo")
	c.Assert(err, IsNil)
	value, err = ctx2.Lookup("foo")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "bar0")

	err = ctx.Unbind("foo")
	c.Assert(err, IsNil)
	value, err = ctx2.Lookup("foo")
	c.Assert(err, IsNil)
	c.Assert(value, IsNil)

	err = ctx.Bind("wombat", "Freddy Boy")
	c.Assert(err, IsNil)
	value, err = ctx2.Lookup("wombat")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "Freddy Boy")
	ctx99 := ctx2.SetParent(nil)
	c.Assert(ctx99, Equals, ctx2)
	c.Assert(ctx2.GetParent(), IsNil)
	value, err = ctx2.Lookup("wombat")
	c.Assert(err, IsNil)
	c.Assert(value, IsNil) // broke chain of contexts
}

func (s *XLSuite) TestSerialization(c *C) {
	var err error
	rng := xr.MakeSimpleRNG()
	n := 16 + rng.Intn(16)
	var keys []string
	var values []string
	mCheck := make(map[string]string)
	for i := 0; i < n; i++ {
		key := rng.NextFileName(8)
		ok := false
		for ok {
			if _, ok = mCheck[key]; !ok {
				break
			}
		}
		// we have a unique key
		val := rng.NextFileName(8)
		mCheck[key] = val
		keys = append(keys, key)
		values = append(values, val)
	}
	c.Assert(len(keys), Equals, n)
	c.Assert(len(values), Equals, n)

	// build a context using these key/value pairs
	ctx := NewNewContext()
	for k, v := range mCheck {
		err = ctx.Bind(k, v)
		c.Assert(err, IsNil)
	}
	ser := ctx.String()
	ss := strings.Split(ser, "\n")
	ss = ss[0 : len(ss)-1] // we split on newline, so last line is empty
	length := len(ss)
	c.Assert(length, Equals, n)
	where := rng.Intn(length)
	indent := rng.Intn(3)
	var spaces string
	for i := 0; i < indent; i++ {
		spaces += " "
	}
	comment := []string{spaces + "# this goes somewhere in the serialization"}
	chunk1 := ss[0:where]
	chunk2 := ss[where:]
	// XXX Apparent bug in Go: if we start with ss2 := chunk1, this fails,
	// with the comment overwriting the last line of chunk1
	var ss2 []string
	ss2 = append(ss2, chunk1...)
	ss2 = append(ss2, comment...)
	ss2 = append(ss2, chunk2...)
	ser = strings.Join(ss2, "\n")
	deser, err := ParseContext(ser)
	c.Assert(err, IsNil)
	c.Assert(deser.Size(), Equals, n)
	for k, v := range mCheck {
		var v2 interface{}
		v2, err = deser.Lookup(k)
		c.Assert(err, IsNil)
		c.Assert(v2, Equals, v)
	}
}
