package om

// xgo/xml/om/elementI.go

type ElementI interface {
	AddNamespace(prefix, namespace string) error
	GetNodeList() *NodeList
	AddChild(elm NodeI) error
	WalkHolders(v VisitorI) error

	GetPrefix() string
	GetName() string
	GetAttrList() *AttrList
	AddAttr(prefix, name, value string) error
	// 	Element AddAttr (name, value string)
	GetAttr(n uint) (*Attr, error)
	WalkAttrs(v VisitorI) error

	NodeI
}
