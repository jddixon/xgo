package om

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestComment(c *C) {
	comment := NewComment("the big guy")
	c.Assert(comment, NotNil)
	c.Assert(comment.GetText(), Equals, "the big guy")
	c.Assert(comment.ToXml(), Equals, "<!-- the big guy -->\n")
}
