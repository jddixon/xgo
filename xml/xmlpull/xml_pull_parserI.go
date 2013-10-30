package xmlpull

import (
	"io"
)

// XmlPull is related to the XML Pull Parser, a Java interface that defines 
// parsing functionlity provided in the XMLPULL V1 API.  Visit 
// http://www.xmlpull.org to learn more about the API and its Java
// implementations.

type XmlPullParserI interface {
    SetFeature(name string, state bool)
    GetFeature(name string) bool
    SetProperty(name string, value interface{})
    GetProperty(name string) interface{}
    SetInput(in io.Reader) 
    DefineEntityReplacementText(entityName, replacementText string)
    GetNamespaceCount(depth int)  int
    GetNamespacePrefix(pos int)  string
    GetNamespaceUri(pos int)  string
    
	// Returns the uri for the given prefix.  The value returned depends
	// upon the current state of the parser.  For example for 'xsi' if the
	// xsi namespace prefix was declared to 'urn:foo' it will return 'urn:foo'.
    // Returns "" if the namespace cannot be found.
    // 
	// This is a convenience method for
    //
	//  for (int i = getNamespaceCount (getDepth ())-1; i >= 0; i--) {
    //   if (getNamespacePrefix (i).equals (prefix)) {
    //     return getNamespaceUri (i);
    //   }
    //  }
    //  return null;
    // However the parser implementation may be more efficient.
	//
	// "ForPrefix" was added to disambiguate.
    GetNamespaceForPrefix (prefix string) string
    
	GetDepth() int
    GetPositionDescription () string
    GetLineNumber() int
    GetColumnNumber() int
    IsWhitespace() bool
    GetText () string
	getTextCharacters(holderForStartAndLength []int) []byte

	// START_TAG / END_TAG SHARED METHODS

    // Returns the namespace URI of the current element.  If namespaces are 
	// NOT enabled, an empty String ("") always is returned.
    // The current event must be START_TAG or END_TAG, otherwise, an error
	// is returned.
    GetNamespace () (string, error)
    
	GetName() string
    GetPrefix() string
    IsEmptyElementTag() bool
    GetAttributeCount() int
    GetAttributeNamespace (index int) string
    GetAttributeName (index int) string
    GetAttributePrefix(index int) string
    
    // Returns the given attribute's value.  The index is zero-based.
    // Returns an IndexOutOfBounds error if the index is out of range
    // or current event type is not START_TAG.
    GetAttributeValue(index int) (string, error)


    // Returns the attributes value identified by the namespace URI and 
	// namespace localName. If namespaces are disabled namespace must be empty.
    // If the current event type is not START_TAG then IndexOutOfBounds
	// error Exception will be returned.
	// 
	// "NS" added to disambiguate.
    GetAttributeValueNS(namespace, name string) (string, error)

	// ACTUAL PARSING METHODS ///////////////////////////////////////

    GetEventType() int
    Next() int
    NextToken() int
    Require (type_ int, namespace, name string)
    ReadText ()  string
}
