package om

// xgo/xml/om/elementI.go

type ElementI interface {
	GetPrefix() string
	GetName() string
	GetAttrList() *AttrList
	AddAttr(prefix, name, value string) error
	// 	Element AddAttr (name, value string)
	GetAttr(n uint) (*Attr, error)
	WalkAttrs(v VisitorI) error
	// IsElement() bool		// in NodeI
	ToString() string
	// ToXml() string		// in NodeI
	HolderI
}
