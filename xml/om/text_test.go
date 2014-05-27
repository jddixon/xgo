package om

import (
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestSingleText(c *C) {
	text := NewText("the big guy")
	c.Assert(text, NotNil)
	c.Assert(text.GetText(), Equals, "the big guy")
	c.Assert(text.ToXml(), Equals, "the big guy")
}
