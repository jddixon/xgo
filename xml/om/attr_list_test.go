package om

import (
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestOneAttrList(c *C) {
	abc := NewNewAttr("abc", "def")
	list := NewAttrList(abc)
	c.Assert(list.Size(), Equals, uint(1))
	gotten, err := list.Get(0)
	c.Assert(err, IsNil)
	c.Assert(gotten, Equals, abc)
	c.Assert(" abc=\"def\"", Equals, list.ToXml())
}

func (s *XLSuite) TestOnePrefixedAttrList(c *C) {
	abc := NewAttr("p", "abc", "def")
	list := NewAttrList(abc)
	c.Assert(list.Size(), Equals, uint(1))
	gotten, err := list.Get(0)
	c.Assert(err, IsNil)
	c.Assert(gotten, Equals, abc)
	c.Assert(" p:abc=\"def\"", Equals, list.ToXml())
}

func (s *XLSuite) TestMultiAttrList(c *C) {
	abc := NewNewAttr("abc", "123")
	def := NewNewAttr("def", "456")
	ghi := NewAttr("z", "ghi", "789")

	list := NewAttrList(abc, def)
	err := list.Add(ghi)
	c.Assert(list.Size(), Equals, uint(3))
	gotten0, err := list.Get(0)
	c.Assert(err, IsNil)
	gotten1, err := list.Get(1)
	c.Assert(err, IsNil)
	gotten2, err := list.Get(2)
	c.Assert(err, IsNil)
	c.Assert(gotten0, Equals, abc)
	c.Assert(gotten1, Equals, def)
	c.Assert(gotten2, Equals, ghi)
	c.Assert(list.ToXml(), Equals, " abc=\"123\" def=\"456\" z:ghi=\"789\"")
}
