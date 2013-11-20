package om

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestXmlDeclParser(c *C) {
	version, encoding, err := ParseXmlDecl(DEFAULT_XML_DECL)
	c.Assert(err, IsNil)
	c.Assert(version, Equals, "1.0")
	c.Assert(encoding, Equals, "UTF-8")
}

func (s *XLSuite) TestEmptyDoc(c *C) {
	doc, err := NewNewDocument()
	c.Assert(err, IsNil)
	c.Assert(doc, NotNil)

	c.Assert(doc.ToXml(), Equals, DEFAULT_XML_DECL)
	c.Assert(doc.GetVersion(), Equals, "1.0")
	c.Assert(doc.GetEncoding(), Equals, "UTF-8")

	elm := doc.GetElementNode()
	c.Assert(elm.GetName(), Equals, "")

	// XXX Ignore for now, but this shows up a problem: because of
	// the way in which Element is included in Document, the line
	//     err = doc.SetDocument(doc2)
	// makes sense.

	//doc2, err := NewNewDocument()
	//c.Assert(err, IsNil)
	//c.Assert(doc2, NotNil)

	//// XXX THIS SHOULD BE NONSENSE
	//err = doc.SetDocument(doc2)
	//c.Assert(err, IsNil)
}

// XXX MUCH MORE TESTING NECESSARY

// XXX UNSATISFACTORY: the order in which namespaces are printed
// is actually unpredictable.
//
// XXX MORE SIGNIFICANT: there is no way to associate namespaces
// with a document.  They are actually associated with the root
// element, which may have unexpected consequences if the document
// is transformed or the root element replaced.
//
func (s *XLSuite) TestAddingNamespaces(c *C) {
	doc, err := NewNewDocument()
	c.Assert(err, IsNil)
	c.Assert(doc, NotNil)

	root, err := NewNewElement("abc")
	c.Assert(err, IsNil)

	root.AddNamespace("c", "http://org.xlattice.xgo/core")
	root.AddNamespace("x", "http://org.xlattice.xgo/xml")

	doc.SetElementNode(root)

	// XXX ORDER OF NAMESPACE DECLARATIONS IS UNPREDICTABLE:
	// have to fiddle with this to get it to succeed
	expected := DEFAULT_XML_DECL + "<abc" +
		" xmlns:c=\"http://org.xlattice.xgo/core\"" +
		" xmlns:x=\"http://org.xlattice.xgo/xml\"" +
		"/>\n"
	c.Assert(doc.ToXml(), Equals, expected)

}
