package om

import (
	"fmt"
)

var _ = fmt.Print // DEBUG

type Node struct {
	doc    DocumentI // this node's ultimate parent; may be nil
	holder ElementI  // this node's immediate parent; may be nil
}

func NewNode() (node *Node) {
	node = &Node{}
	return
}

// Return this node's ultimate parent, the XML document.
func (node *Node) GetDocument() DocumentI {
	return node.doc
}

// Change this node's ultimate parent, its Document.
//
// This gets overridden in Holder and possibly elsewhere.
//
// XXX Warning: possibility of introducing cycles or inconsistencies.
//
// @param newDoc new value assigned; may be nil
//
func (node *Node) SetDocument(newDoc DocumentI) (err error) {
	if newDoc == nil {
		err = NilDocument
	} else {
		node.doc = newDoc
	}
	return
}

// Get this node's parent.
//
func (node *Node) GetHolder() ElementI {
	return node.holder
}

// Set or change this node's immediate parent; also change
// this node's ultimate parent if necessary.
//
// XXX There is no check for inconsistencies with the
// parent's NodeList, nor against the introduction of cycles
// into the node graph.
//
// @param h a reference to the new parent, may be nil
//
func (node *Node) SetHolder(e ElementI) {
	if e == nil {
		node.doc = nil
	} else {
		eDoc := e.GetDocument()
		if eDoc == nil {
			node.doc = nil
		} else if eDoc != node.doc {
			node.SetDocument(eDoc)
		}
	}
}

// VISITOR-RELATED///////////////////////////////////////////////
//

// Walk a Visitor through a Node.  This is overridden when
// suitable by subclasses.

func (node *Node) WalkAll(v VisitorI) (err error) {
	err = v.OnEntry(node)

	// Holders also visit their NodeLists
	if err == nil {
		err = v.OnExit(node)
	}
	return
}

// // EVAL /////////////////////////////////////////////////////////
// public final Boolean evalAsBoolean(String s) {
//     // STUB
//     return Boolean.FALSE
// }
// //
//  * XXX XPath refers to this as evalAsLocation?
//
// public final NodeSet evalAsNodeSet(String s) {
//     // STUB
//     return nil
// }
// public final Numeric evalAsNumeric(String s) {
//     // STUB
//     return new Numeric(0.0)
// }
// public final String evalAsString(String s) {
//     // STUB
//     return nil
// }
// // XXX SHOULD BE NodeSet? XXX
// public final Boolean evalAsBoolean(Context ctx, String s) {
//     // STUB
//     return Boolean.FALSE
// }
// public final Node evalAsLocation(Context ctx, String s) {
//     // STUB
//     return nil
// }
// public final Numeric evalAsNumeric(Context ctx, String s) {
//     // STUB
//     return new Numeric(0.0)
// }
// public final String evalAsString(Context ctx, String s) {
//     // STUB
//     return nil
// } // GEEP

// TYPE IDENTIFIERS /////////////////////////////////////////////

// one of these gets overridden in each subclass
func (node *Node) IsAttr() bool { return false }

// one of these gets overridden in each subclass
func (node *Node) IsCdata() bool { return false }

// one of these gets overridden in each subclass
func (node *Node) IsComment() bool { return false }

// one of these gets overridden in each subclass
func (node *Node) IsDocument() bool { return false }

// one of these gets overridden in each subclass
func (node *Node) IsDocType() bool { return false }

// one of these gets overridden in each subclass
func (node *Node) IsElement() bool { return false }

// CDATA subclass of Text
// one of these gets overridden in each subclass
func (node *Node) IsText() bool { return false }

// one of these gets overridden in each subclass
func (node *Node) IsPI() bool {
	return false
}

// SERIALIZATION ////////////////////////////////////////////////
// this should be changed to follow the same pattern as expr,
// supporting indenting
func (node *Node) ToXml() (s string) {
	// XXX ABSTRACT
	return
}
