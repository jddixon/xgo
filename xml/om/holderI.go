package om

// xgo/xml/om/holderI.go

type HolderI interface {
	AddNamespace(prefix, namespace string) error
	GetNodeList() *NodeList
	// SetDocument(newDoc *Document) error	// in NodeI
	AddChild(elm NodeI) error
	// WalkAll(v VisitorI) error			// in NodeI
	WalkHolders(v VisitorI) error
	NodeI
}
