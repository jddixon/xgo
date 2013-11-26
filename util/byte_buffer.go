package util

import ()

// A ByteBuffer is something that satisfies io.Writer

type ByteBuffer struct {
	buf []byte
}

func NewByteBuffer(size int) (b *ByteBuffer, err error) {
	if size <= 0 {
		err = NonPositiveBufferSize
	} else {
		buf := make([]byte, 0, size)
		b = &ByteBuffer{
			buf: buf,
		}
	}
	return
}

func (b *ByteBuffer) Cap() int {
	return cap(b.buf)
}
func (b *ByteBuffer) Len() int {
	return len(b.buf)
}
func (b *ByteBuffer) Write(p []byte) (n int, err error) {
	// simple-minded implementation
	b.buf = append(b.buf, p...)
	n = len(p)
	return
}

func (b *ByteBuffer) String() string {
	return string(b.buf)
}
