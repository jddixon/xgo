package util

// xgo/util/byte_buffer_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestByteBuffer(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_BYTE_BUFFER")
	}

	var soFar, totalSize int

	rng := xr.MakeSimpleRNG()
	count := 8 + rng.Intn(8)
	p := make([][]byte, count)  // we make this many little slices
	sizes := make([]int, count) // the length of each
	for i := 0; i < count; i++ {
		size := 16 + rng.Intn(16)
		sizes[i] = size
		p[i] = make([]byte, size)
		rng.NextBytes(&p[i]) // fill the slice with random values
		totalSize = totalSize + size
	}
	capacity := 2 * totalSize
	b, err := NewByteBuffer(capacity)
	c.Assert(err, IsNil)
	c.Assert(b.Len(), Equals, 0)
	c.Assert(b.Cap(), Equals, capacity)

	for i := 0; i < count; i++ {
		n, err := b.Write(p[i])
		c.Assert(err, IsNil)
		c.Assert(n, Equals, sizes[i])
		soFar += n
		c.Assert(b.Len(), Equals, soFar)
	}
}
