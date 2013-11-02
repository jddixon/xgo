package om

// xgo/xml/om/doc_setter.go


// Walk down the subtrees, child nodes and their descendents,
// setting each node's Document to match this Holder's
// Document.
//
// Implements VisitorI
//
type DocSetter struct {
	ThisDocument	*Document
}

func NewDocSetter(h *Holder) (ds *DocSetter, err error) {
	if n == nil {
		err = NilHolder
	} else {
		ds = &DocSetter{
			ThisDocument:	h.GetDocument(),
		}
	}
	return
}
		
// On arriving at the node, set its Document.
func (ds *DocSetter) OnEntry (n *Node ) (err error) {
	if n == nil {
		err = NilNode
	} else {
		n.SetDocument(ds.ThisDocument)
	}
}
// On leaving, do nothing.//
func (ds *DocSetter) OnExit (n *Node) (err error) { 
	if n == nil {
		err = NilNode
	}
	return
}
