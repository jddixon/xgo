package om

import (
	"strings"
)

// A container for attributes.  The order of attributes in the
// container is not significant and is not guaranteed to be
// repeatable.
type AttrList struct {
	attrs  []*Attr
	holder ElementI // points back to the element
}

// Create an empty attribute list with no holder specified.
func NewNewAttrList() *AttrList {
	return &AttrList{}
}

// Create an attribute list from an arbitrary number of attributes.
func NewAttrList(a ...*Attr) (aList *AttrList) {
	aList = NewNewAttrList()
	for i := 0; i < len(a); i++ {
		aList.attrs = append(aList.attrs, a[i])
	}
	return
}

// Add an attribute to an existing container.
func (aList *AttrList) Add(attr *Attr) (err error) {
	if attr == nil {
		err = NilAttr
	} else {
		if aList.holder != nil {
			attr.SetHolder(aList.holder)
		}
		aList.attrs = append(aList.attrs, attr)
	}
	return
}

// Insert an attribute into an existing container in a particular
// place, displacing any existing attributes if necessary.
//
// @param n    zero-based index at which the Attr is to be inserted
// @param attr the attribute to be inserted
// @throws IndexOutOfBoundsException if n is negative or out of range
// @throws NullPointerException if the Attr argument is null
//
func (aList *AttrList) Insert(n uint, attr *Attr) (err error) {
	if attr == nil {
		err = NilAttr
	} else if n > uint(len(aList.attrs)) {
		err = IndexOutOfBounds
	}
	if err == nil {
		attr.SetHolder(aList.holder)
		if n == uint(len(aList.attrs)) {
			err = aList.Add(attr)
		} else {
			head := aList.attrs[:n]
			tail := aList.attrs[n:]
			aList.attrs = append(head, attr)
			aList.attrs = append(aList.attrs, tail...)
		}
	}
	return
}

// Get the Nth attribute.
//
// @param n index of the Attr to be returned
// Return the Nth attr in the list
// @throws IndexOutOfBoundsException
//
func (aList *AttrList) Get(n uint) (attr *Attr, err error) {
	if n >= uint(len(aList.attrs)) {
		err = IndexOutOfBounds
	} else {
		attr = aList.attrs[n]
	}
	return
}

// Return number of attrs in the list.
func (aList *AttrList) Size() uint {
	return uint(len(aList.attrs))
}

// PROPERTIES ///////////////////////////////////////////////////

// Return the Element that the attribute belongs to//
func (aList *AttrList) GetHolder() HolderI {
	return aList.holder
}

//
// Set the Holder for this attribute.  By definition the Holder
// must be an XML Element.
//
// @param h the Holder being assigned
//
func (aList *AttrList) SetHolder(elm HolderI) {
	e := elm.(*Element)
	aList.holder = e
	for i := uint(0); i < aList.Size(); i++ {
		aList.attrs[i].SetHolder(e)
	}
}

// VISITOR-RELATED///////////////////////////////////////////////

// Walk a Visitor through the list of attributes, visiting each
// in turn.
// @param v the visitor
//
func (aList *AttrList) WalkAll(v VisitorI) (err error) {
	for i := uint(0); i < aList.Size(); i++ {
		err = aList.attrs[i].WalkAll(v)
		if err != nil {
			break
		}
	}
	return
}

// SERIALIZATION ////////////////////////////////////////////////

// Return the list in XML String form.
//
func (aList *AttrList) ToXml() string {

	var ss []string
	for i := uint(0); i < aList.Size(); i++ {
		ss = append(ss, aList.attrs[i].ToXml())
	}
	return strings.Join(ss, "")
}
