package om

import (
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestAttrWithoutPrefix(c *C) {
	attr := NewNewAttr("georgieBoy", "the big guy")
	c.Assert(attr, NotNil)
	c.Assert("", Equals, attr.GetPrefix())
	c.Assert("georgieBoy", Equals, attr.GetName())
	c.Assert("the big guy", Equals, attr.GetValue())
	c.Assert(" georgieBoy=\"the big guy\"", Equals, attr.ToXml())
}
func (s *XLSuite) TestPrefixedAttrs(c *C) {
	attr := NewAttr("a", "b", "c")
	c.Assert("a", Equals, attr.GetPrefix())
	c.Assert(" a:b=\"c\"", Equals, attr.ToXml())
}
