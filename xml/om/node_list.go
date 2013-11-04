package om

// xgo/xml/om/node_list.go

import ()

const DEFAULT_NODELIST_SIZE = 8 // XXX DROP ASAP

// A container for Nodes.  Each Holder (Document or Element) has a
// NodeList, but the reverse is not necessarily true.
// <p/>
// XXX NEED TO SIMPLIFY THE HANDLING OF THE ELEMENT NODE IN A
// XXX DOCUMENT
type NodeList struct {

	// list of child nodes//
	nodes []*Node
	// immediate parent, might be nil//
	holder *Holder
	// ultimate parent, might be nil*/
	doc *Document
}

// Create an empty node list.
func NewNewNodeList(node *Node) *NodeList {
	var nodes []*Node
	return &NodeList{
		nodes: nodes,
	}
}

// Create a node list with only one member.
func NewNodeList(node *Node) *NodeList {
	nodes := []*Node{node}
	return &NodeList{
		nodes: nodes,
	}
}

// Add a Node to the NodeList.
//
// XXX Should check for cycles; if the Holder is a document,
// XXX there may be only one Element node.
//
// @param node the node to be appended
// @return a reference to this list, to allow chaining
// @throws NullPointerException if the Node argument is nil
//
func (nl *NodeList) Append(node *Node) (this *NodeList) {
	this = nl
	node.setHolder(this)
	nl.nodes = append(nl.nodes, node)
	return this
}

// Copy the nodes from another NodeList into this one, then
// delete them from the source, to ease GC.
//
// @throws NullPointerException if otherList is nil
//
func (nl *NodeList) MoveFrom(otherList *NodeList) (this *NodeList, err error) {
	if otherList == nil {
		err = EmptyOtherList
	} else {
		this = nl
		for i := uint(0); i < otherList.Size(); i++ {
			var node *Node
			node, err = otherList.Get(i)
			if err != nil {
				break
			}
			node.SetHolder(this)
			nl.nodes = append(nl.nodes(node))
		}
	}
	if err == nil {
		otherList.Clear()
	}
	return
}
func (nl *NodeList) Clear() {
	var nodes []*Node
	nl.nodes = nodes
}

//
// @param n    zero-based index at which the Node is to be inserted
// @param node the node to be inserted
// @throws IndexOutOfBoundsException if n is negative or out of range
// @throws NullPointerException if the Node argument is nil
//
func (nl *NodeList) Insert(n uint, node *Node) (err error) {
	if n > nl.Size() {
		err = IndexOutOfBounds
	}
	if err == nil && node == nil {
		err = NilNode
	}
	if err == nil {
		node.SetHolder(nl)
		if n == nl.Size() {
			nl.nodes = append(nl.nodes, node)
		} else {
			head := nl.nodes[0:n]
			tail := nl.nodes[n:]
			nl.nodes = append(head, node)
			nl.nodes = append(nl.nodes, tail...)
		}
	}
	return
}

// Return whether there are no nodes in the list
func (nl *NodeList) IsEmpty() bool {
	return len(nl.nodes) == 0
}

// Return the Nth node in the list.
//
// @param n index of the Node to be returned
// @return the Nth node in the list
// @throws IndexOutOfBoundsException
//
func (nl *NodeList) Get(n uint) (node *Node, err error) {
	if n >= nl.Size() {
		err = IndexOutOfBounds
	} else {
		node = nl.nodes[n]
	}
	return
}

// Return number of nodes in the list
//
func (nl *NodeList) Size() uint {
	return uint(len(nl.nodes))
}

// PROPERTIES ///////////////////////////////////////////////////
// @return the immediate parent of this list//
func (nl *NodeList) GetHolder() *Holder {
	return nl.holder
}

//
// Change the immediate parent of this list, here and in
// descendent nodes.
//
// XXX SHOULD CHECK FOR GRAPH CYCLES
//
// @param h the new parent; may be nil
//
func (nl *NodeList) SetHolder(h *Holder) {
	var doc *Document
	if h == nil {
		doc = nil
	} else {
		doc = h.GetDocument()
	}
	for i := uint(0); i < nl.Size(); i++ {
		nl.Get(i).SetHolder(h)
	}
}

// VISITOR-RELATED///////////////////////////////////////////////
// Take the visitor through every node in the list, recursing.//
func (nl *NodeList) WalkAll(v VisitorI) (err error) {
	for i := uint(0); i < nl.Size(); i++ {
		err = nl.Get(i).WalkAll(v)
		if err != nil {
			break
		}
	}
	return
}

//
// Take the Visitor through the list, visiting any node which is
// a Holder, recursively.  Used when you don't want to visit, for
// example, attributes.
//
func (nl *NodeList) WalkHolders(v VisitorI) (err error) {
	for i := uint(0); err == nil && i < nl.Size(); i++ {
		var n *Node
		n, err = nl.Get(i)
		if err == nil && n.IsHolder() {
			err = n.WalkHolders(v)
		}
	}
	return
}

// SERIALIZATION METHODS ////////////////////////////////////////
//
// A String containing each of the Nodes in XML form, recursively,
// without indenting.
//
func (nl *NodeList) ToXml() (s string) {
	for i := uint(0); i < nl.Size(); i++ {
		var node *Node
		node, _ = nl.Get(i)
		s += node.ToXml()
	}
	return
}
