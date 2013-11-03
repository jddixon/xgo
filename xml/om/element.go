package om

import (
	"fmt"
	"strings"
)

// An XML element node.  This implementation specifies the element
// in terms of its prefix and name.
//
// XXX Need to consider whether to store the Namespace (uri) rather
// than the prefix; could resolve prefix in the constructor.  This
// decision will have consequences.
//
// In this implementation an element always has its own attribute
// list.
//
type Element struct {
	prefix string
	name   string
	aList  *AttrList
	Holder
}

//Create an XML element, given its prefix and name.  Both
//prefix and name should conformant to the XML specifications
//and must not contain colons (that is, they must be NCNames).
//
//@param prefix NCName or nil
//@param name   NCName, must not be nil
//
func NewElement(prefix, name string) (e *Element, err error) {
	// super()

	aList, err := NewAttrList()
	if err == nil {
		e = &Element{
			prefix: prefix,
			name:   name,
			aList:  aList,
		}
		aList.SetHolder(e)
	}
	return
}

//
//Create an XML element, defaulting the prefix to nil.
//
func (e *Element) NewNewElement(name string) (*Element, error) {
	return NewElement("", name)
}

// PROPERTIES ///////////////////////////////////////////////////
// @return the prefix, an NCName ""
func (e *Element) GetPrefix() string {
	return e.Prefix
}

// @return the element name, an NCName, which may not be nil//
func (e *Element) GetName() string {
	return e.Name
}

//@return the attribute list - may be empty, may not be nil
//
func (e *Element) GetAttrList() *AttrList {
	return e.aList
}

// ATTRIBUTES ///////////////////////////////////////////////////

//Add an attribute to this element.
//@param prefix to attribute name, may be nil
//@param name   the attribute name itself
//@param value  the String value the attribute is set to
//@return       a reference to this Element, to allow chaining
//
func (e *Element) AddAttr(prefix, name, value string) (
	this *Element, err error) {

	this = e
	attr := NewAttr(prefix, name, value)
	_, err = e.aList.Add(attr)
	return
}

//Add an element, defaulting its prefix to nil.
//
// func (e *Element) Element AddAttr (name, value string) {
//    return addAttr (nil, name, value)
// }

//@param  n index of the parameter to be returned
//@return the Nth attribute
//
func (e *Element) GetAttr(n uint) (*Attr, error) {
	return e.aList.Get(n)
}

// VISITOR-RELATED///////////////////////////////////////////////

func (e *Element) WalkAttrs(v *Visitor) error {
	return e.aList.WalkAll(v)
}

// NODE METHODS /////////////////////////////////////////////////

func (e *Element) IsElement() bool {
	return true
}

//Preliminary version, for debugging.
//
//@return the element in string form, without its attributes
//

func (e *Element) ToString() string {
	return fmt.Sprintf("[Element: tag: %s ...]", e.Name)
}

//@return the element and its attributes in XML form, unindented
//
func (e *Element) ToXml() (s string) {

	// conditionally output prefix

	s = "<" + e.Name

	// conditionally output attributes
	attrCount := e.aList.Size()
	for i = uint(0); i < attrCount; i++ {
		attr, _ := e.aList.Get(i)
		s += " " + attr.ToXml()
	}

	// conditionally output ns2pf
	for i := 0; i < len(e.nsUris); i++ {
		ns := nsUris[i]
		p := e.ns2pf[ns]
		s += " "
		if p == "" {
			s += "xmlns=\""
		} else {
			s += "xmlns:" + p + "=\""
		}
		s += ns + "\""
	}

	if len(e.nodes) > 0 {
		// line separator
		s += ">\n"
		ss := []string{s}

		// conditionally output body
		for i := 0; i < len(e.nodes); i++ {
			body := e.nodes[i].ToXml()
			ss = append(ss, body)
		}
		// prefix ?
		ss = append(ss, "</"+e.name+">\n")
		s += strings.Join(ss, "\n")

	} else {
		// empty element
		s += "/>\n"
	}
	return
}
