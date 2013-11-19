package om

// xgo/xml/om/node_list_test.go

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestOneNodeList(c *C) {
	elm, err := NewNewElement("abc")
	c.Assert(err, IsNil)
	list := NewNodeList(elm)
	c.Assert(list.Size(), Equals, uint(1))
	list0, err := list.Get(0)
	c.Assert(err, IsNil)
	c.Assert(list0, Equals, elm)
}

func (s *XLSuite) TestCopyingNodeLists(c *C) {
	tom, err := NewNewElement("tom")
	c.Assert(err, IsNil)
	dick, err := NewNewElement("dick")
	c.Assert(err, IsNil)
	harry, err := NewNewElement("harry")
	c.Assert(err, IsNil)
	joe, err := NewNewElement("joe")
	c.Assert(err, IsNil)
	elms := []*Element{tom, dick, harry, joe}
	list := NewNewNodeList()
	list2 := NewNewNodeList()
	for i := 0; i < len(elms); i++ {
		list.Append(elms[i])
	}
	list2.MoveFrom(list)
	c.Assert(list.Size(), Equals, uint(0))
	c.Assert(list2.Size(), Equals, uint(len(elms)))
	for i := uint(0); i < uint(len(elms)); i++ {
		elm, err := list2.Get(i)
		c.Assert(err, IsNil)
		c.Assert(elm, Equals, elms[i])
	}
	list2.Clear()
	c.Assert(list2.Size(), Equals, uint(0))
}
