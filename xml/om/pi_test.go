package om

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestSimplePI(c *C) {
	pi := NewProcessingInstruction("perl", "chomp;")
	c.Assert(pi, NotNil)
	c.Assert(pi.GetTarget(), Equals, "perl")
	c.Assert(pi.GetText(), Equals, "chomp;")
	c.Assert(pi.ToXml(), Equals, "<?perl chomp;?>\n")
}

func (s *XLSuite) TestComboPI(c *C) {
	pi, err := ProcessingInstructionFromString("perl chomp;")
	c.Assert(err, IsNil)
	c.Assert(pi, NotNil)
	c.Assert(pi.GetTarget(), Equals, "perl")
	c.Assert(pi.GetText(), Equals, "chomp;")
	c.Assert(pi.ToXml(), Equals, "<?perl chomp;?>\n")
}
