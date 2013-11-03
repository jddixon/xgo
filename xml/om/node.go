package om

// With apologies to org.xlattice.corexml.om.Node.java

import ()

type Node struct {
	// this node's ultimate parent; may be null
	doc *Document
	// this node's immediate parent; may be null
	holder *Holder
}

func NewNode() (node *Node) {
	node = &Node{}
	return
}

// Return this node's ultimate parent, the XML document.
func (node *Node) GetDocument() *Document {
	return node.doc
}

// Change this node's ultimate parent, its Document.
//
// This gets overridden in Holder and possibly elsewhere.
//
//  XXX Warning: possibility of introducing cycles or inconsistencies.
//
//  @param newDoc new value assigned; may be null
//
func (node *Node) SetDocument(newDoc *Document) {
	node.doc = newDoc
}

// Get this node's parent.
//
func (node *Node) GetHolder() *Holder {
	return node.holder
}

// Set or change this node's immediate parent; also change
// this node's ultimate parent if necessary.
//
// XXX There is no check for inconsistencies with the
// parent's NodeList, nor against the introduction of cycles
// into the node graph.
//
// @param h a reference to the new parent, may be null
//
func (node *Node) SetHolder(h *Holder) {
	if h == nil {
		node.doc = nil
	} else if h.getDocument() != node.doc {
		node.SetDocument(h.GetDocument())
	}
}

// VISITOR-RELATED///////////////////////////////////////////////
//
// Runs down the graph counting nodes by type; another Visitor.
//
//public class NodeCounter implements Visitor {
//    private int attrCount
//    private int commentCount
//    private int docCount
//    private int docTypeCount
//    private int elementCount
//    private int piCount
//    private int textCount
//
//    public NodeCounter () { }
//
//    public void onEntry (Node n) {
//        if (n.isAttr())                     attrCount++
//        if (n.isComment())                  commentCount++
//        if (n.isDocument())                 docCount++
//        if (n.isDocType())                  docTypeCount++
//        if (n.isElement())                  elementCount++
//        if (n.isProcessingInstruction())    piCount++
//        if (n.isText())                     textCount++
//    }
//    public void onExit (Node n) { }
//    // PROPERTIES /////////////////////////////////////
//    public int attrCount()      { return attrCount; }
//    public int commentCount()   { return commentCount; }
//    public int docCount()       { return docCount; }
//    public int docTypeCount()   { return docTypeCount; }
//    public int elementCount()   { return elementCount; }
//    public int piCount()        { return piCount; }
//    public int textCount()      { return textCount; }
//}
//
//
// Walk a Visitor through a Node.  This is overridden when
// suitable by subclasses.

//public void walkAll (Visitor v) {
//    v.onEntry(this)
//    // Holders also visit their NodeLists
//    v.onExit(this)
//}
//  // EVAL /////////////////////////////////////////////////////////
//  public final Boolean evalAsBoolean(String s) {
//      // STUB
//      return Boolean.FALSE
//  }
//  //
//   * XXX XPath refers to this as evalAsLocation?
//
//  public final NodeSet evalAsNodeSet(String s) {
//      // STUB
//      return nil
//  }
//  public final Numeric evalAsNumeric(String s) {
//      // STUB
//      return new Numeric(0.0)
//  }
//  public final String evalAsString(String s) {
//      // STUB
//      return nil
//  }
//  // XXX SHOULD BE NodeSet? XXX
//  public final Boolean evalAsBoolean(Context ctx, String s) {
//      // STUB
//      return Boolean.FALSE
//  }
//  public final Node evalAsLocation(Context ctx, String s) {
//      // STUB
//      return nil
//  }
//  public final Numeric evalAsNumeric(Context ctx, String s) {
//      // STUB
//      return new Numeric(0.0)
//  }
//  public final String evalAsString(Context ctx, String s) {
//      // STUB
//      return nil
//  }

// TYPE IDENTIFIERS /////////////////////////////////////////////
// one of these gets overridden in each subclass
func (node *Node) IsAttr() bool { return false }

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
func (node *Node) IsProcessingInstruction() bool {
	return false
}

// SERIALIZATION ////////////////////////////////////////////////
// this should be changed to follow the same pattern as expr,
// supporting indenting
func (node *Node) ToXml() string {
	// XXX ABSTRACT
	return
}
