package om

import (
	"fmt"
	"strings"
)

// An XML element node.  This implementation specifies the element
// in terms of its prefix and name.
//
// XXX Need to consider whether to store the Namespace (uri) rather
// than the prefix; could resolve prefix in the constructor.  This
// decision will have consequences.
//
// In this implementation an element always has its own attribute
// list.
//
type Element struct {
	prefix   string
	name     string
	aList    *AttrList
	nodeList *NodeList
	nsUris   []string

	// maps namespaces into prefixes
	ns2pf map[string]string
	// reverse mapping, prefixes into namespaces
	pf2ns map[string]string

	Node
}

//Create an XML element, given its prefix and name.  Both
//prefix and name should conformant to the XML specifications
//and must not contain colons (that is, they must be NCNames).
//
//@param prefix NCName or nil
//@param name   NCName, must not be nil
//
func NewElement(prefix, name string) (e *Element, err error) {

	aList := NewAttrList()

	var nsUris []string
	if err == nil {
		e = &Element{
			prefix:   prefix,
			name:     name,
			aList:    aList,
			ns2pf:    make(map[string]string), // namespace --> prefix map
			pf2ns:    make(map[string]string), // prefix --> namespace map
			nsUris:   nsUris,
			nodeList: NewNewNodeList(), // creates empty list
		}
		aList.SetHolder(e)
		e.nodeList.SetHolder(e)
	}
	return
}

//
//Create an XML element, defaulting the prefix to nil.
//
func NewNewElement(name string) (*Element, error) {
	return NewElement("", name)
}

// PROPERTIES ///////////////////////////////////////////////////
// @return the prefix, an NCName ""
func (e *Element) GetPrefix() string {
	return e.prefix
}

// @return the element name, an NCName, which may not be nil//
func (e *Element) GetName() string {
	return e.name
}

//@return the attribute list - may be empty, may not be nil
//
func (e *Element) GetAttrList() *AttrList {
	return e.aList
}

// Add a prefix-namespace pair, updating the maps.
//
// @param prefix    the prefix, a NCNAME, may not be null
// @param namespace XML-compatible namespace
//
func (e *Element) AddNamespace(prefix, namespace string) (err error) {
	// XXX NEED MORE REASONABLE CHECKS
	if namespace == "" {
		err = EmptyNamespace
	}
	e.ns2pf[namespace] = prefix
	if prefix != "" {
		e.pf2ns[prefix] = namespace
	}
	e.nsUris = append(e.nsUris, namespace) // SILLY LEVEL OF OVERKILL
	return
}

// Return a pointer to the list of children of this Element.
// XXX This is not secure.
func (e *Element) GetNodeList() *NodeList {
	return e.nodeList
}

// Set this Element's ultimate parent, the Document it belongs
// to.
//
func (e *Element) SetDocument(newDoc DocumentI) (err error) {
	var docSetter *DocSetter

	e.doc = newDoc
	docSetter, err = NewDocSetter(e) // will use e's Document

	if err == nil {
		err = e.GetAttrList().WalkAll(docSetter)
	}
	if err == nil {
		e.nodeList.WalkAll(docSetter) // set in subtree
	}
	return
}

// ATTRIBUTES ///////////////////////////////////////////////////

//Add an attribute to this element.
//@param prefix to attribute name, may be nil
//@param name   the attribute name itself
//@param value  the String value the attribute is set to
//@return       a reference to this Element, to allow chaining
//
func (e *Element) AddAttr(prefix, name, value string) (err error) {

	attr := NewAttr(prefix, name, value)
	err = e.aList.Add(attr)
	return
}

//Add an element, defaulting its prefix to nil.
//
// func (e *Element) Element AddAttr (name, value string) {
//    return addAttr (nil, name, value)
// }

//@param  n index of the parameter to be returned
//@return the Nth attribute
//
func (e *Element) GetAttr(n uint) (*Attr, error) {
	return e.aList.Get(n)
}

// CHILDREN /////////////////////////////////////////////////////////
//
// Add a child Node to the Element.
//
func (e *Element) AddChild(elm NodeI) (err error) {
	if elm == nil {
		err = NilChild
	} else {
		e.nodeList.nodes = append(e.nodeList.nodes, elm)
	}
	return
}

// VISITOR-RELATED///////////////////////////////////////////////////

func (e *Element) WalkAttrs(v VisitorI) error {
	return e.aList.WalkAll(v)
}

// Take a Visitor on that walk down the subtrees, visiting
// every Node.
//
func (e *Element) WalkAll(v VisitorI) (err error) {
	err = v.OnEntry(e)
	if err == nil {
		err = e.WalkAttrs(v)
	}
	if err == nil {
		err = e.nodeList.WalkAll(v)
		if err == nil {
			err = v.OnExit(e)
		}
	}
	return
}

// Take a Visitor on that walk down the subtrees, visiting
// only subnode which are themselves Holders.
//
func (e *Element) WalkHolders(v VisitorI) (err error) {
	err = v.OnEntry(e)
	if err == nil {
		err = e.nodeList.WalkHolders(v)
		if err == nil {
			err = v.OnExit(e)
		}
	}
	return
}

// NODE METHODS /////////////////////////////////////////////////

func (e *Element) IsElement() bool {
	return true
}

//Preliminary version, for debugging.
//
//@return the element in string form, without its attributes
//

func (e *Element) ToString() string {
	return fmt.Sprintf("[Element: tag: %s ...]", e.name)
}

// Return the element and its attributes in XML form, unindented.
// If the element has no name, we return an empty string.
//
func (e *Element) ToXml() (s string) {

	// conditionally output prefix
	if e.name != "" {
		s = "<" + e.name

		// conditionally output attributes
		if e.aList != nil {
			attrCount := e.aList.Size()
			for i := uint(0); i < attrCount; i++ {
				attr, _ := e.aList.Get(i)
				s += " " + attr.ToXml()
			}
		}

		// conditionally output ns2pf
		for i := 0; i < len(e.nsUris); i++ {
			ns := e.nsUris[i]
			p := e.ns2pf[ns]
			s += " "
			if p == "" {
				s += "xmlns=\""
			} else {
				s += "xmlns:" + p + "=\""
			}
			s += ns + "\""
		}

		nodes := e.GetNodeList()
		if (nodes != nil) && (nodes.Size() > 0) {
			// line separator
			s += ">\n"
			ss := []string{s}

			// conditionally output body
			for i := uint(0); i < nodes.Size(); i++ {
				node, _ := nodes.Get(i)
				body := node.ToXml()
				ss = append(ss, body)
			}
			// prefix ?
			ss = append(ss, "</"+e.name+">\n")
			s += strings.Join(ss, "\n")

		} else {
			// empty element
			s += "/>\n"
		}
	}
	return
}

// OTHER METHODS ////////////////////////////////////////////////
//
// XXX Replace h *Holder with e *Element throughout.
//
// Arrive here having seen either START_DOCUMENT or START_TAG and
// having created a Node of the appropriate type.
//
// throws CoreXmlException, IOException, XmlPullParserException

// func (h *Holder) Populator (xpp pp.XmlPullParserI, depth, endEvent int) (
// 	err error) {
//
//     if (!nodeList.isEmpty())
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
//                 nodeList.append(elm)
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
//                 nodeList.append( new Text(xpp.getText()))
//                 break
//             case pp.COMMENT:
//                 nodeList.append( new Comment(xpp.getText()))
//                 break
//             case pp.CDSECT:
//                 nodeList.append( new Cdata(xpp.getText()))
//                 break
//             case pp.PROCESSING_INSTRUCTION:
//                 nodeList.append( new ProcessingInstruction (xpp.getText() ))
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
