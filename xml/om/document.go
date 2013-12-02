package om

// xgo/xml/om/document.go

import (
	"fmt"
	"strings"
)

var _ = fmt.Print

const (
	DEFAULT_VERSION  = "1.0"
	DEFAULT_ENCODING = "UTF-8"
	DEFAULT_XML_DECL = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
)

// A simplistic XML decl parser.
//
func ParseXmlDecl(decl string) (version, encoding string, err error) {

	var ndx int
	decl = strings.TrimSpace(decl)
	if !strings.HasPrefix(decl, "<?xml ") {
		err = IllFormedDecl
	} else {
		decl = decl[6:]
		ndx = strings.Index(decl, "?>")
		if ndx == -1 {
			err = IllFormedDecl
		} else {
			decl = decl[:ndx]
		}
	}
	if err == nil {
		parts := strings.Split(decl, " ")
		if len(parts) == 2 {
			subParts := strings.Split(parts[0], "\"")
			if len(subParts) != 3 {
				err = IllFormedDecl
			} else if subParts[0] == "version=" {
				version = subParts[1]
			} else {
				err = IllFormedDecl
			}
			if err == nil {
				subParts = strings.Split(parts[1], "\"")
				if len(subParts) != 3 {
					err = IllFormedDecl
				} else if subParts[0] == "encoding=" {
					encoding = subParts[1]
				} else {
					err = IllFormedDecl
				}
			}

		} else {
			err = IllFormedDecl
		}
	}
	return
}

// An XML Document, a Holder which can contain only one Element instead
// of a NodeList, has no attributes, and no namespaces.
type Document struct {
	version  string
	encoding string
	docType  *DocumentType
	Element
}

// Create an XML document with the XML declaration passed.
//
// @param decl the XML declaration
//
func NewDocumentFromDecl(decl string) (doc *Document, err error) {

	version, encoding, err := ParseXmlDecl(decl)
	if err == nil {
		doc = &Document{
			version:  version,
			encoding: encoding,
		}
	}
	return
}

// Create an XML document with the default XML declaration.
//
func NewNewDocument() (*Document, error) {
	return NewDocumentFromDecl(DEFAULT_XML_DECL)
}

// Create an XML document with the version number and encoding
// specified.
//
// XXX CHECKS NEEDED
//
// @param version  XML version number; if nil, uses the default
// @param encoding if nil, the default is used
//
func NewDocument(version, encoding string) (doc *Document, err error) {

	if version == "" {
		version = DEFAULT_VERSION
	}
	if encoding == "" {
		encoding = DEFAULT_ENCODING
	}
	doc = &Document{
		version:  version,
		encoding: encoding,
	}
	return
}

// PROPERTIES ///////////////////////////////////////////////////

// @return the XML document type
//
func (doc *Document) GetDocType() *DocumentType {
	return doc.docType
}

func (doc *Document) SetDocType(dt *DocumentType) (err error) {
	if dt == nil {
		err = NilDocType
	} else {
		doc.docType = dt
	}
	return
}

// @return the XML encoding used in the document//
func (doc *Document) GetEncoding() string {
	return doc.encoding
}

// Get the document's element node; there may only be one.
//
// @return a reference to the document's element node
//
func (doc *Document) GetElementNode() ElementI {
	return &doc.Element
}

// Set the document's element node.   There may only be one such element.
//
// XXX There must be some checks to ensure that the
// element is well-formed AND that this does not introduce
// cycles into the graph.
//
func (doc *Document) SetElementNode(elm ElementI) (err error) {

	if elm == nil {
		err = NilElement
	} else {
		ePtr := elm.(*Element)
		ePtr.SetHolder(doc)
		ePtr.SetDocument(doc)
		doc.Element = *ePtr
	}
	return
}

// @return the XML version of this document//
func (doc *Document) GetVersion() string {
	return doc.version
}

// NODE METHODS /////////////////////////////////////////////////

// Return true; this node is an Document.
func (a *Document) IsDocument() bool {
	return true
}

// Generate the XML document in String form.  The standard XML
// declaration is prefixed.  This method traverses the entire
// document recursively.  The document is <b>not</b> indented.
//
// @return the entire document in String form
//
func (doc *Document) ToXml() (s string) {
	s = "<?xml version=\"" + doc.version +
		"\" encoding=\"" + doc.encoding + "\"?>\n"
	s += doc.Element.ToXml()
	return
}
