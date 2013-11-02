package om

import (
	"fmt"
)

//XML attribute, either a triplet (prefix, name, value) or a pair
//(name, value).  The prefixed form is a 'colonized' string like
//"prefix:name", where <i>colonized</i> means 'containing a colon'.
//As required by the XML standard, neither <code>prefix</code> nor
//<code>name</code> may contain a colon.  In the actual XML, the
//value will be doubly (") or singly (') quoted.
//
type Attr struct {
	Prefix string
	Name   string
	Value  string
	Node
}

//Create an attribute.
//
//@param prefix NCNAME (non-colonized name) identifying the namespace
//@param name   attribute name, another NCNAME
//@param value  the attribute's value
//
func NewAttr(prefix, name, value string) (a *Attr) {
	a = &Attr{
		Prefix: prefix,
		Name:   name,
		Value:  value,
	}
	return a
}

//Default constructor with nil prefix.
//
func NewNewAttr(name, value string) {
	return NewAttr("", name, value)
}

// PROPERTIES ///////////////////////////////////////////////////

// return the prefix part of the name; may be nil//
func (a *Attr) GetPrefix() {
	return a.Prefix
}

// @return the unqualified name of the attribute//
func (a *Attr) GetName() {
	return a.Name
}

// @return the value assigned to the attribute//
func (a *Attr) GetValue() {
	return a.Value
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
// @return true; this node is an Attr//
func (a *Attr) IsAttr() bool {
	return true
}

//Convert the Node to XML form.  If the prefix is nil, it is
//omitted.
//
//@return the attribute in XML form
//
func (a *Attr) ToXml() (s string) {

	if a.Prefix != "" {
		s = fmt.Sprintf(" %s:%s=\"%s\"", a.Prefix, a.Name, a.Value)
	} else {
		s = fmt.Sprintf(" %s=\"%s\"")
	}
	return
}
