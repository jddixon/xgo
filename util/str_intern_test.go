package util

// xgo/util/str_intern_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
	"reflect"
	"unsafe"
)

func (s *XLSuite) TestStrIntern(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_STR_INTERN")
	}

	x := "<section>abcdef0123456789abcdef0123456789"
	y := x
	// the compiler isn't smart enough to see that this is the same string
	z := "<section>"
	z += "abcdef0123456789abcdef0123456789"

	xPtr := unsafe.Pointer(&x)
	yPtr := unsafe.Pointer(&y)
	zPtr := unsafe.Pointer(&z)

	// reflect.StringHeader has Data uintptr and Len int fields
	xHdr := (*reflect.StringHeader)(xPtr)
	yHdr := (*reflect.StringHeader)(yPtr)
	zHdr := (*reflect.StringHeader)(zPtr)

	c.Assert(xHdr.Data, Equals, yHdr.Data)
	c.Assert(xHdr.Data, Not(Equals), zHdr.Data)
	c.Assert(yHdr.Data, Not(Equals), zHdr.Data)

	// Intern them ...
	si := NewStrIntern()
	x = si.Intern(x)
	y = si.Intern(y)
	z = si.Intern(z)

	xPtr = unsafe.Pointer(&x)
	yPtr = unsafe.Pointer(&y)
	zPtr = unsafe.Pointer(&z)

	xHdr = (*reflect.StringHeader)(xPtr)
	yHdr = (*reflect.StringHeader)(yPtr)
	zHdr = (*reflect.StringHeader)(zPtr)

	// x, y, and z should be conventionally equal to one another
	// but should also have identical uintptrs

	c.Assert(x, Equals, y)
	c.Assert(y, Equals, z)
	c.Assert(x, Equals, x)

	c.Assert(xHdr.Data, Equals, yHdr.Data)
	c.Assert(xHdr.Data, Equals, zHdr.Data)
}
