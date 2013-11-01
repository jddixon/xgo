package xmlpull

import (
	"io"
)

/////////////////////////////////////////////////////////////////////
// XXX COMMENTS to be merged into this file; as text is merged in here,
// it should be deleted from COMMENTS
/////////////////////////////////////////////////////////////////////

// XmlPull is related to the XML Pull Parser, a Java interface that defines
// parsing functionlity provided in the XMLPULL V1 API.  Visit
// http://www.xmlpull.org to learn more about the API and its Java
// implementations.

type XmlPullParserI interface {
	// May return UnsupportedFeature or EmptyFeatureString
	SetFeature(name string, state bool) error

	// May return EmptyFeatureString
	GetFeature(name string) (bool, error)

	// May return XmlPullParserException
	SetProperty(name string, value interface{}) error

	GetProperty(name string) interface{}

	// May return XmlPullParserException
	SetInput(in io.Reader) error

	// May return XmlPullParserException
	DefineEntityReplacementText(entityName, replacementText string) error

	// May return XmlPullParserException
	GetNamespaceCount(depth int) (int, error)

	// May return XmlPullParserException
	GetNamespacePrefix(pos int) (string, error)

	// May return PosOutOfRange
	GetNamespaceUri(pos int) (string, error)

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
	// May return XmlPullParserException
	//
	// "ForPrefix" was added to disambiguate.
	GetNamespaceForPrefix(prefix string) string

	// MISCELLANEOUS REPORTING METHODS //////////////////////////////

	GetDepth() int

	GetPositionDescription() string
	GetLineNumber() int
	GetColumnNumber() int

	// TEXT-RELATED METHODS

	// May return NotElementContentRelated
	IsWhitespace() (bool, error)

	GetText() string
	getTextCharacters(holderForStartAndLength []int) []byte

	// START_TAG / END_TAG SHARED METHODS ///////////////////////////

	// Returns the namespace URI of the current element.  If namespaces are
	// NOT enabled, an empty String ("") always is returned.
	// The current event must be START_TAG or END_TAG, otherwise, an error
	// is returned.
	GetNamespace() (string, error)

	GetName() string
	GetPrefix() string

	// May return NotOnStartTag
	IsEmptyElementTag() (bool, error)

	// START_TAG ATTRIBUTE RETRIEVAL METHODS ////////////////////////

	GetAttributeCount() int
	GetAttributeNamespace(index int) string
	GetAttributeName(index int) string

	// May throw IndexOutOfBounds
	GetAttributePrefix(index int) string

	// Returns the given attribute's value.  The index is zero-based.
	// Returns an IndexOutOfBounds error if the index is out of range
	// or current event type is not START_TAG.
	GetAttributeValue(index int) (string, error)

	// Returns the attributes value identified by the namespace URI and
	// namespace localName. If namespaces are disabled namespace must be empty.
	// If the current event type is not START_TAG then IndexOutOfBounds
	// error will be returned.
	//
	// "NS" added to disambiguate.
	GetAttributeValueNS(namespace, name string) (string, error)

	// ACTUAL PARSING METHODS ///////////////////////////////////////

	// May throw XmlPullParserException
	GetEventType() (int, error)

	// May throw XmlPullParserException or IOError
	Next() (int, error)

	// May throw XmlPullParserException or IOError
	NextToken() (int, error)

	// UTILITY METHODS //////////////////////////////////////////////

	// May throw XmlPullParserException or IOError
	Require(type_ int, namespace, name string)

	// May throw XmlPullParserException or IOError
	ReadText() string
}
