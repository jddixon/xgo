package om

// xgo/xml/om/node_counter.go

import ()

// Runs down the graph counting nodes by type; another Visitor.
//
// Implements VisitorI
type NodeCounter struct {
	attrCount    int
	commentCount int
	docCount     int
	docTypeCount int
	elementCount int
	piCount      int
	textCount    int
}

func NewNodeCounter() *NodeCounter {
	return &NodeCounter{}
}

func (nc *NodeCounter) OnEntry(n *Node) {
	if n.IsAttr() {
		nc.attrCount++
	}
	if n.IsComment() {
		nc.commentCount++
	}
	if n.IsDocument() {
		nc.docCount++
	}
	if n.IsDocType() {
		nc.docTypeCount++
	}
	if n.IsElement() {
		nc.elementCount++
	}
	if n.IsProcessingInstruction() {
		nc.piCount++
	}
	if n.IsText() {
		nc.textCount++
	}
}
func (nc *NodeCounter) OnExit(n *Node) {
}

// PROPERTIES /////////////////////////////////////
func (nc *NodeCounter) AttrCount() int {
	return nc.attrCount
}
func (nc *NodeCounter) CommentCount() int {
	return nc.commentCount
}
func (nc *NodeCounter) DocCount() int {
	return nc.docCount
}
func (nc *NodeCounter) DocTypeCount() int {
	return nc.docTypeCount
}
func (nc *NodeCounter) ElementCount() int {
	return nc.elementCount
}
func (nc *NodeCounter) PiCount() int {
	return nc.piCount
}
func (nc *NodeCounter) TextCount() int {
	return nc.textCount
}
