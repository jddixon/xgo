package om

// xgo/xml/om/holder.go

import (
// pp "github.com/jddixon/xgo/xml/xmlpull"
)

// A Holder is something which can have children and namespaces, so a
// Document or an Element.
type Holder struct {
	// every Holder has a list of child Nodes//
	nodes  *NodeList
	nsUris []string

	// maps namespaces into prefixes//
	ns2pf map[string]string
	// reverse mapping, prefixes into namespaces//
	pf2ns map[string]string

	Node
}

//
// Create a Holder, associated data structures, and a pair of
// Visitors.
//
func NewHolder() (h *Holder) {
	// super ()
	// nodes.setHolder(this)

	// SILLY OVERKILL //
	ns2pf := make(map[string]string) // namespace --> prefix map
	pf2ns := make(map[string]string) // prefix --> namespace map
	var nsUris []string

	h = &Holder{
		ns2pf: ns2pf,
		pf2ns: pf2ns,
	}
	return
}

// Add a prefix-namespace pair, updating the maps.
//
// XXX SHOULD BE MOVED TO Element XXX
// @param prefix    the prefix, a NCNAME, may not be null
// @param namespace XML-compatible namespace
//
func (h *Holder) AddNamespace(prefix, namespace string) (err error) {
	// XXX NEED MORE REASONABLE CHECKS
	if namespace == "" {
		err = EmptyNamespace
	}
	h.ns2pf[namespace] = prefix
	if prefix != "" {
		h.pf2ns[prefix] = namespace
	}
	h.nsUris = append(h.nsUris, namespace) // SILLY LEVEL OF OVERKILL
	return
}

// PROPERTIES ///////////////////////////////////////////////////

// Return a pointer to the list of children of this Holder.
// XXX This is not secure.
func (h *Holder) GetNodeList() *NodeList {
	return h.nodes
}

// Set this Holder's ultimate parent, the Document it belongs
// to.
//
func (h *Holder) SetDocument(newDoc *Document) (err error) {
	var docSetter *DocSetter
	if h.IsDocument() {
		err = SettingDocsDoc
	} else {
		h.doc = newDoc
		docSetter, err = NewDocSetter(h) // will use h's Document
	}
	if err == nil {
		if h.IsElement() {
			err = h.GetAttrList().WalkAll(docSetter)
		}
	}
	if err == nil {
		h.nodes.WalkAll(docSetter) // set in subtree
	}
	return
}

// OTHER METHODS ////////////////////////////////////////////////
//
// Add a child Node to the Holder.
//
// @param elm  child Node to be added
// @throws     NullPointerException if the child is nil
//
func (h *Holder) AddChild(elm *Node) (err error) {
	if elm == nil {
		err = NilChild
	} else {
		h.nodes = append(h.nodes, elm)
	}
	return
}

// VISITOR-RELATED///////////////////////////////////////////////
//
// Take a Visitor on that walk down the subtrees, visiting
// every Node.
//
func (h *Holder) WalkAll(v VisitorI) (err error) {
	err = v.OnEntry(h)
	if err == nil && h.IsElement() {
		err = h.WalkAttrs(v)
	}
	if err == nil {
		err = h.nodes.WalkAll(v)
		if err == nil {
			err = v.OnExit(h)
		}
	}
	return
}

// Take a Visitor on that walk down the subtrees, visiting
// only subnodes which are themselves Holders.
//
func (h *Holder) WalkHolders(v VisitorI) (err error) {
	h = v.OnEntry(h)
	if err == nil {
		err = h.nodes.walkHolders(v)
		if err == nil {
			h = v.OnExit(h)
		}
	}
	return
}

// OTHER METHODS ////////////////////////////////////////////////
//
// Arrive here having seen either START_DOCUMENT or START_TAG and
// having created a Node of the appropriate type.
//
// throws CoreXmlException, IOException, XmlPullParserException

// func (h *Holder) Populator (xpp pp.XmlPullParserI, depth, endEvent int) (
// 	err error) {
//
//     if (!nodes.isEmpty())
//         throw new IllegalStateException("NodeList is not empty")
//     int elementCount = 0
//     int event
//
//     // COLLECT ANY NAME SPACES ////////////////////////
//     int myDepth = xpp.getDepth()
//     int nsPrev  = myDepth <= 0 ? 0 : xpp.getNamespaceCount(myDepth - 1)
//     int nsNow   = xpp.getNamespaceCount(myDepth)
//     nsUris = new ArrayList (nsNow - nsPrev);  // XXX CHECK ME
//
//     for (int i = nsPrev; i < nsNow; i++) {
//         String prefix = xpp.getNamespacePrefix(i)
//         String uri    = xpp.getNamespaceUri(i)
//         addNamespace (prefix, uri)
//         // DEBUG
//         //System.out.println("namespace " + i + ", " + prefix + ":" + uri)
//         // END
//     }
//     // COLLECT ATTRIBUTES /////////////////////////////
//     if (isElement()) {
//         int count = xpp.getAttributeCount()
//         Element me = (Element)this
//         for (int i = 0; i < count; i++) {
//             // IGNORE TYPE FOR NOW
//             // IGNORE ATTR NAMESPACE
//             me.addAttr(xpp.getAttributePrefix(i),
//                         xpp.getAttributeName(i), xpp.getAttributeValue(i))
//         }
//     }
//     // COLLECT CHILDREN ///////////////////////////////
//     // detect empty document
//     try {
//         event = xpp.nextToken()
//     } catch (IOException ioe) {
//         return
//     }
// 		// empty document detection did nextToken()
//     for (event != pp.END_DOCUMENT && event != endEvent
//                     event = xpp.nextToken()) {
//         switch (event) {
//             case pp.START_TAG:
//                 if (isDocument() && elementCount > 0)
//                     throw new CoreXmlException(
//                             "more than one root element found")
//                 elementCount++
//                 Element elm = new Element(xpp.getName())
//                 elm.populator(xpp, depth + 1, pp.END_TAG)
//                 nodes.append(elm)
//                 if (isDocument()) {
//                     Document me = (Document) this
//                     if (me.getElementNode() == nil)
//                         me.setElementNode(elm)
//                     else
//                         throw new IllegalStateException (
//                             "second element at root level in document")
//                 }
//                 break
//             case pp.IGNORABLE_WHITESPACE:
//             case pp.TEXT:
//                 nodes.append( new Text(xpp.getText()))
//                 break
//             case pp.COMMENT:
//                 nodes.append( new Comment(xpp.getText()))
//                 break
//             case pp.CDSECT:
//                 nodes.append( new Cdata(xpp.getText()))
//                 break
//             case pp.PROCESSING_INSTRUCTION:
//                 nodes.append( new ProcessingInstruction (xpp.getText() ))
//                 break
//
//             // //////////////////////////////////////////////////
//             // THESE ARE NOT YET HANDLED ////////////////////////
//             // //////////////////////////////////////////////////
//             case pp.DOCDECL:
//                 // DEBUG
//                 System.out.println("    *** IGNORING DOCDECL TOKEN ***")
//                 // END
//                 break
//             case pp.ENTITY_REF:
//                 // DEBUG
//                 System.out.println("    *** IGNORING ENTITY_REF TOKEN ***")
//                 // END
//                 break
//             default:
//                 throw new CoreXmlException(
//                     "unknown event type " + event)
//         }
//     }
// }
//
