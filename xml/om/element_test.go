package om

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestEmptyElement(c *C) {
	element, err := NewNewElement("abc")
	c.Assert(err, IsNil)
	c.Assert(element, NotNil)
	c.Assert(element.GetPrefix(), Equals, "")
	c.Assert(element.GetName(), Equals, "abc")
	c.Assert(element.ToXml(), Equals, "<abc/>\n")
}

// Test added 2006-06-21 after deserializing problem discovered.
// Unfortunately, the problem doesn't show up here.
//
var data060621 = DEFAULT_XML_DECL +
	"<project>" +
	"  <!-- project description -->" +
	"  <description>" +
	"    Dummy parent project." +
	"  </description>" +
	"  <shortDescription>" +
	"    dummy parent project" +
	"  </shortDescription>" +
	"</project>"

	// XXX This won't work without an XML pull parser
	//
	//func (s *XLSuite) Test21060621 (c *C) {
	//    Document doc = new XmlParser(new StringReader(data060621)).read()
	//    String serialization = doc.toXml()
	//    // DEBUG
	//    System.out.println(serialization)
	//    // END
	//    assertSameSerialization( data060621, serialization )
	//}
	// XXX TEST ATTRIBUTE LIST **********

	// XXX MUCH MORE TESTING NEEDED *****
