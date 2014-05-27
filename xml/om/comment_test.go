package om

import (
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestComment(c *C) {
	comment := NewComment("the big guy")
	c.Assert(comment, NotNil)
	c.Assert(comment.GetText(), Equals, "the big guy")
	c.Assert(comment.ToXml(), Equals, "<!-- the big guy -->\n")
}
