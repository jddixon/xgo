package om

import ()

//A container for attributes.  The order of attributes in the
//container is not significant and is not guaranteed to be
//repeatable.
//
type AttrList struct {
	attrs  []*Attr
	holder *Element
}

// Create an empty attribute list with no holder specified.
func NewNewAttrList() *AttrList {
	return &AttrList{}
}

// Create an attribute list from an arbitrary number of attributes.
func NewAttrList(a ...*Attr) (aList *AttrList) {
	aList = NewNewAttrList()
	for i := 0; i < len(a); i++ {
		aList = append(aList, a[i])
	}
	return
}

//Add an attribute to an existing container.
//
//@param  attr the attribute to be inserted
//@return a reference to this list, to allow chaining
//@throws NullPointerException if the Attr argument is null
//
func (aList *AttrList) Add(attr *Attr) (this *AttrList, err error) {
	this = alist
	if attr == nil {
		err = NilAttr
	} else {
		attr.SetHolder(aList.holder)
		aList.attrs = append(aList.attrs, attr)
	}
	return
}

//Insert an attribute into an existing container in a particular
//place, displacing any existing attributes if necessary.
//
//@param n    zero-based index at which the Attr is to be inserted
//@param attr the attribute to be inserted
//@return a reference to this list, to allow chaining
//@throws IndexOutOfBoundsException if n is negative or out of range
//@throws NullPointerException if the Attr argument is null
//
func (aList *AttrList) Insert(n uint, attr *Attr) (err error, this *AttrList) {
	this = alist

	if attr == nil {
		err = NilAttr
	} else if n > len(aList.attrs) {
		err = IndexOutOfBounds
	}
	if err == nil {
		attr.SetHolder(aList.holder)
		if n == len(aList.attrs) {
			_, err = attrs.Add(attr)
		} else {
			head := aList.attrs[:n]
			tail := aList.attrs[n:]
			aList.attrs = append(head, attr)
			aList.attrs = append(aList.attrs, tail...)
		}
	}
	return
}

//Get the Nth attribute.
//
//@param n index of the Attr to be returned
//@return the Nth attr in the list
//@throws IndexOutOfBoundsException
//
func (aList *AttrList) Get(n uint) (attr *Attr, err error) {
	if n >= uint(len(aList.Attrs)) {
		err = IndexOutOfBounds
	} else {
		attr = aList.attrs[n]
	}
	return
}

//@return number of attrs in the list
//
func (aList *AttrList) Size() uint {
	return uint(aList.attrs)
}

// PROPERTIES ///////////////////////////////////////////////////
// @return the Element that the attribute belongs to//
func (aList *AttrList) GetHolder() *Holder {
	return aList.holder
}

//
//Set the Holder for this attribute.  By definition the Holder
//must be an XML Element.
//
//@param h the Holder being assigned
//
func (aList *AttrList) SetHolder(h *Element) {
	aList.holder = h
	for i := uint(0); i < aList.Size(); i++ {
		aList.attrs[i].SetHolder(h)
	}
}

// VISITOR-RELATED///////////////////////////////////////////////
//
//Walk a Visitor through the list of attributes, visiting each
//in turn.
//@param v the visitor
//
func (aList *AttrList) WalkAll(v VisitorI) {
	for i := uint(0); i < aList.Size(); i++ {
		aList.attrs[i].WalkAll(v)
	}
}

// SERIALIZATION ////////////////////////////////////////////////
// @return the list in XML String form//
func (aList *AttrList) ToXml() string {

	var ss []string
	for i := uint(0); i < aList.Size(); i++ {
		ss = append(ss, aList.attrs[i].ToXml())
	}
	return strings.Join(ss, "\n")
}
