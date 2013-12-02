package om

import (
	"fmt"
)

// XML attribute, either a triplet (prefix, name, value) or a pair
// (name, value).  The prefixed form is a 'colonized' string like
// "prefix:name", where <i>colonized</i> means 'containing a colon'.
// As required by the XML standard, neither <code>prefix</code> nor
// <code>name</code> may contain a colon.  In the actual XML, the
// value will be doubly (") or singly (') quoted.
//
type Attr struct {
	prefix string
	name   string
	value  string
	Node
}

/////////////////////////////////////////////////////////////////////
// XXX If I comment out everything below this block, I still get
// "attr.go:18: invalid recursive type Node
// followed by
// "attr.go:19: invalid recursive type Attr
// on the lines above.
/////////////////////////////////////////////////////////////////////

// Create an attribute.
//
// @param prefix NCNAME (non-colonized name) identifying the namespace
// @param name   attribute name, another NCNAME
// @param value  the attribute's value
//
func NewAttr(prefix, name, value string) (a *Attr) {
	a = &Attr{
		prefix: prefix,
		name:   name,
		value:  value,
	}
	return
}

// Default constructor with nil prefix.
//
func NewNewAttr(name, value string) *Attr {
	return NewAttr("", name, value)
}

// PROPERTIES ///////////////////////////////////////////////////

// return the prefix part of the name; may be nil//
func (a *Attr) GetPrefix() string {
	return a.prefix
}

// @return the unqualified name of the attribute//
func (a *Attr) GetName() string {
	return a.name
}

// @return the value assigned to the attribute//
func (a *Attr) GetValue() string {
	return a.value
}

// VISITOR-RELATED///////////////////////////////////////////////
//
//Method used by classes walking the XML document.
//
//@param v the Visitor walking the document
//
func (a *Attr) WalkAll(v VisitorI) (err error) {
	err = v.OnEntry(a)
	if err == nil {
		err = v.OnExit(a)
	}
	return
}

// NODE METHODS /////////////////////////////////////////////////

// Return true; this node is an Attr.
func (a *Attr) IsAttr() bool {
	return true
}

//Convert the Node to XML form.  If the prefix is nil, it is
//omitted.
//
//@return the attribute in XML form
//
func (a *Attr) ToXml() (s string) {

	if a.prefix != "" {
		s = fmt.Sprintf(" %s:%s=\"%s\"", a.prefix, a.name, a.value)
	} else {
		s = fmt.Sprintf(" %s=\"%s\"", a.name, a.value)
	}
	return
}
