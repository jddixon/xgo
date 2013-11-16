package om

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestSimpleCDATA(c *C) {
	cdata := NewCdata("the big guy")
	c.Assert(cdata, NotNil)
	c.Assert(cdata.GetText(), Equals, "the big guy")
	c.Assert(cdata.ToXml(), Equals, "<![CDATA[the big guy]]>\n")
}
