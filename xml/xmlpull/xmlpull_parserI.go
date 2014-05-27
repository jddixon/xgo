package xmlpull

// xgo/xml/xmlpull/xmlpull_parserI.go

import (
	"io"
)

// XmlPull is derived from the XML Pull Parser, a Java interface that defines
// parsing functionlity provided in the XMLPULL V1 API.  Visit
// http://www.xmlpull.org to learn more about the API and its Java
// implementations.
//
// 2013-11-07: The comments that follow are  taken directly from the XML Pull
// Parser spec with little modification.  As time permits these will be
// edited to suit the Go environment.
//
// <p>There are following different
// kinds of parser depending on which features are set:<ul>
// <li>behaves like XML 1.0 comliant non-validating parser
//  <em>if no DOCDECL is present</em> in XML documents when
//   FEATURE_PROCESS_DOCDECL is false (this is <b>default parser</b>
//   and internal enetites can still be defiend with defineEntityReplacementText())
// <li>non-validating parser as defined in XML 1.0 spec when
//   FEATURE_PROCESS_DOCDECL is true
// <li>validating parser as defined in XML 1.0 spec when
//   FEATURE_VALIDATION is true (and that implies that FEATURE_PROCESS_DOCDECL is true)
// </ul>
//
// <p>There are only two key methods: next() and nextToken() that provides
// access to high level parsing events and to lower level tokens.
//
// <p>The parser is always in some event state and type of the current event
// can be determined by calling
// <a href="#next()">getEventType()</a>.
// Initially the parser is in the
// <a href="#START_DOCUMENT">START_DOCUMENT</a> state.
//
// <p>Method <a href="#next()">next()</a> returns int that contains identifier of parsing event.
// This method can return following events (and will change parser state to the returned event):<dl>
// <dt><a href="#START_TAG">START_TAG</a><dd> XML start tag was read
// <dt><a href="#TEXT">TEXT</a><dd> element contents was read and is available via getText()
// <dt><a href="#END_TAG">END_TAG</a><dd> XML end tag was read
// <dt><a href="#END_DOCUMENT">END_DOCUMENT</a><dd> no more events is available
// </dl>
//
// BEING HACKED, hence the odd mixture of Go and Java syntax.
//
// A minimal working example of the use of the API would look like this:
//
// <pre>
// import (
//     "fmt"a
//     "io"
//     "strings"
// )
//
// func main() {
//
//     var rd1 io.Reader = string.NewReader("<foo>Hello, my good world!</foo>" )
//     xpp, _ := NewNewParser(rd1)	// accept default encoding, UTF-8
//     curEvent := xpp.getCurEvent()
//     while curEvent != xpp.END_DOCUMENT {
//         if curEvent == xpp.START_DOCUMENT {
//             fmt.Println("Start document")
//         } else if curEvent == xpp.END_DOCUMENT {
//             fmt.Printfln ("End document")
//         } else if curEvent == xpp.START_TAG {
//             fmt.Printf ("Start tag %s\n", xpp.getName())
//         } else if curEvent == xpp.END_TAG {
//             fmt.Printf ("End tag %s\n, "+xpp.getName())
//         } else if curEvent == xpp.TEXT {
//             fmt.Printf ("Text %s\n", xpp.getText())
//         }
//         curEvent = xpp.Next()
//     }
// }
// </pre>
//
// <p>When run it will produce the following output:
// <pre>
// Start document
// Start tag foo
// Text Hello, my good world!
// End tag foo
// </pre>
//
// <p>For more details on use of API please read
// Quick Introduction available at <a href="http://www.xmlpull.org">http://www.xmlpull.org</a>
//
// @see XmlPullParserFactory
// @see #defineEntityReplacementText
// @see #next
// @see #nextToken
// @see #FEATURE_PROCESS_DOCDECL
// @see #FEATURE_VALIDATION
// @see #START_DOCUMENT
// @see #START_TAG
// @see #TEXT
// @see #END_TAG
// @see #END_DOCUMENT
//
type XmlPullParserI interface {

	// -- SetFeature ------------------------------------------------

	// Use this call to change the general behaviour of the parser,
	// such as namespace processing or doctype declaration handling.
	// This method must be called before the first call to next or
	// nextToken. Otherwise, an exception is trown.
	// <p>Example: call setFeature(FEATURE_PROCESS_NAMESPACES, true) in order
	// to switch on namespace processing. Default settings correspond
	// to properties requested from the XML Pull Parser factory
	// (if none were requested then all feautures are by default false).
	//
	// @exception XmlPullParserException if feature is not supported or can not be set
	// @exception IllegalArgumentException if feature string is null
	//
	// May return UnsupportedFeature or EmptyFeatureString
	SetFeature(name string, state bool) error

	// -- GetFeature ------------------------------------------------

	// Return the current value of the feature with given name.
	// NOTE: unknown features are <string>always</strong> returned as false
	//
	// @param name The name of feature to be retrieved.
	// @return The value of named feature.
	// @exception IllegalArgumentException if feature string is null
	//
	// May return EmptyFeatureString
	//
	GetFeature(name string) (bool, error)

	// -- SetProperty -----------------------------------------------

	// Set the value of a property.
	//
	// The property name is any fully-qualified URI.
	//
	// May return XmlPullParserException
	//
	SetProperty(name string, value interface{}) error

	// -- GetProperty -----------------------------------------------

	// Look up the value of a property.
	//
	// The property name is any fully-qualified URI. I
	// NOTE: unknown features are <string>always</strong> returned as null
	//
	// @param name The name of property to be retrieved.
	// @return The value of named property.
	//
	GetProperty(name string) interface{}

	// -- SetInput --------------------------------------------------

	// Set the input for parser. Parser event state is set to START_DOCUMENT.
	// Using null parameter will stop parsing and reset parser state
	// allowing parser to free internal resources (such as parsing buffers).
	//
	// May return XmlPullParserException
	SetInput(in *io.Reader) error

	// -- DefineEntityReplacement -----------------------------------

	// Set new value for entity replacement text as defined in
	// <a href="http://www.w3.org/TR/REC-xml#intern-replacement">XML 1.0 Section 4.5
	// Construction of Internal Entity Replacement Text</a>.
	// If FEATURE_PROCESS_DOCDECL or FEATURE_VALIDATION are set then calling this
	// function will reulst in exception because when processing of DOCDECL is enabled
	// there is no need to set manually entity replacement text.
	//
	// <p>The motivation for this function is to allow very small implementations of XMLPULL
	// that will work in J2ME environments and though may not be able to process DOCDECL
	// but still can be made to work with predefined DTDs by using this function to
	// define well known in advance entities.
	// Additionally as XML Schemas are replacing DTDs by allowing parsers not to process DTDs
	// it is possible to create more efficient parser implementations
	// that can be used as underlying layer to do XML schemas validation.
	//
	//
	// <p><b>NOTE:</b> this is replacement text and it is not allowed
	//  to contain any other entity reference
	// <p><b>NOTE:</b> list of pre-defined entites will always contain standard XML
	// entities (such as &amp;amp; &amp;lt; &amp;gt; &amp;quot; &amp;apos;)
	// and they cannot be replaced!
	//
	// @see #setInput
	// @see #FEATURE_PROCESS_DOCDECL
	// @see #FEATURE_VALIDATION
	//
	// May return XmlPullParserException
	DefineEntityReplacementText(entityName, replacementText string) error

	// -- GetNamespaceCount -----------------------------------------

	// Return position in stack of first namespace slot for element at passed depth.
	// If namespaces are not enabled it returns always 0.
	// <p><b>NOTE:</b> default namespace is not included in namespace table but
	//  available by getNamespace() and not available from getNamespace(String)
	//
	// @see #getNamespacePrefix
	// @see #getNamespaceUri
	// @see #getNamespace()
	// @see #getNamespace(String)
	//
	// May return XmlPullParserException
	//
	GetNamespaceCount(elmDepth int) (int, error)

	// -- GetNamespacePrefix ----------------------------------------

	// Return namespace prefixes for position pos in namespace stack
	//
	// May return XmlPullParserException
	GetNamespacePrefix(pos int) (string, error)

	// -- GetNamespaceUri -------------------------------------------

	// Return namespace URIs for position pos in namespace stack
	// If pos is out of range it throw exception.
	//
	// May return PosOutOfRange
	GetNamespaceUri(pos int) (string, error)

	// -- GetNamespaceForPrefix -------------------------------------

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
	GetNamespaceForPrefix(prefix string) (string, error)

	// MISCELLANEOUS REPORTING METHODS //////////////////////////////

	// -- GetDepth --------------------------------------------------

	// Returns the current depth of the element.
	// Outside the root element, the depth is 0. The
	// depth is incremented by 1 when a start tag is reached.
	// The depth is decremented AFTER the end tag
	// event was observed.
	//
	// <pre>
	// &lt;!-- outside --&gt;     0
	// &lt;root>               1
	//   sometext           1
	//     &lt;foobar&gt;         2
	//     &lt;/foobar&gt;        2
	// &lt;/root&gt;              1
	// &lt;!-- outside --&gt;     0
	// &lt;/pre&gt;
	// </pre>
	//
	GetDepth() int

	// -- GetPositionDescription ------------------------------------

	// Short text describing parser position, including a
	// description of the current event and data source if known
	// and if possible what parser was seeing lastly in input.
	// This method is especially useful to give more meaningful error messages.
	//
	GetPositionDescription() string

	// -- GetLineNumber ---------------------------------------------

	// Current line number: numebering starts from 1.
	//
	GetLineNumber() int

	// -- GetColumnNumber -------------------------------------------

	// Current column: numbering starts from 0 (returned when parser is
	// in START_DOCUMENT state!)
	//
	GetColumnNumber() int

	// TEXT-RELATED METHODS =========================================

	// -- IsWhiteSpace ----------------------------------------------

	// Check if current TEXT event contains only whitespace characters.
	// For IGNORABLE_WHITESPACE, this is always true.
	// For TEXT and CDSECT if the current event text contains at least one
	// non white space character then false is returned. For any other event
	// an error is returned.
	//
	// NOTE:  non-validating parsers are not
	// able to distinguish whitespace and ignorable whitespace
	// except from whitespace outside the root element. ignorable
	// whitespace is reported as separate event which is exposed
	// via nextToken only.
	//
	// NOTE: this function can be only called for element content related events
	// such as TEXT, CDSECT or IGNORABLE_WHITESPACE otherwise
	// exception will be thrown!
	//
	// May return NotElementContentRelated
	IsWhitespace() (bool, error)

	// -- GetText ---------------------------------------------------

	// Read text content of the current event as String.
	//
	GetText() string

	// -- GetTextCharacters -----------------------------------------

	// Get the buffer that contains text of the current event and
	// start offset of text is passed in first slot of input int array
	// and its length is in second slot.
	//
	// NOTE: this buffer must not be modified and its content MAY change
	// after call to next() or nextToken().
	//
	// NOTE: this methid must return always the same value as getText()
	// and if getText() returns null then this methid returns null as well and
	// values returned in holder MUST be -1 (both start and length).
	//
	// @see #getText
	//
	// @param holderForStartAndLength the 2-element int array into which
	//   values of start offset and length will be written into first and
	//   second slot of array.
	// @return char buffer that contains text of current event
	//  or null if the current event has no text associated.
	//
	GetTextCharacters(holderForStartAndLength []int) []byte

	// START_TAG / END_TAG SHARED METHODS ===========================

	// -- GetNamespace ----------------------------------------------

	// Returns the namespace URI of the current element.  If namespaces are
	// NOT enabled, an empty String ("") always is returned.
	// The current event must be START_TAG or END_TAG, otherwise, an error
	// is returned.
	GetNamespace() (string, error)

	// --- GetName --------------------------------------------------

	// Returns the (local) name of the current element
	// when namespaces are enabled
	// or raw name when namespaces are disabled.
	// The current event must be START_TAG or END_TAG, otherwise "" is returned.
	// NOTE: to reconstruct raw element name
	//  when namespaces are enabled you will need to
	//  add prefix and colon to localName if prefix is not empty.
	//
	GetName() string

	// -- GetPrefix -------------------------------------------------

	// Returns the prefix of the current element
	// or null if elemet has no prefix (is in defualt namespace).
	// If namespaces are not enabled it always returns "".
	// If the current event is not START_TAG or END_TAG "" is returned.
	//
	GetPrefix() string

	// -- IsEmptyElement --------------------------------------------

	// Returns true if the current event is START_TAG and the tag is
	// degenerate (e.g. &lt;foobar/&gt;).
	// NOTE: if parser is not on START_TAG then will return NotOnStartTag
	//
	IsEmptyElementTag() (bool, error)

	// START_TAG ATTRIBUTE RETRIEVAL METHODS ========================

	// -- GetAttributeCount -----------------------------------------

	// Returns the number of attributes on the current element;
	// -1 if the current event is not START_TAG
	//
	// @see #getAttributeNamespace
	// @see #getAttributeName
	// @see #getAttributePrefix
	// @see #getAttributeValue
	//
	GetAttributeCount() int

	// -- GetAttributeNamespace -------------------------------------

	// Returns the namespace URI of the specified attribute
	//  number index (starts from 0).
	// Returns empty string ("") if namespaces are not enabled or attribute has no namespace.
	// Throws an IndexOutOfBoundsException if the index is out of range
	// or current event type is not START_TAG.
	//
	// <p><strong>NOTE:</p> if FEATURE_REPORT_NAMESPACE_ATTRIBUTES is set
	// then namespace attributes (xmlns:ns='...') amust be reported
	// with namespace
	// <a href="http://www.w3.org/2000/xmlns/">http://www.w3.org/2000/xmlns/</a>
	// (visit this URL for description!).
	// The default namespace attribute (xmlns="...") will be reported with empty namespace.
	// Then xml prefix is bound as defined in
	// <a href="http://www.w3.org/TR/REC-xml-names/#ns-using">Namespaces in XML</a>
	// specification to "http://www.w3.org/XML/1998/namespace".
	//
	// @param zero based index of attribute
	// @return attribute namespace or "" if namesapces processing is not enabled.
	//
	GetAttributeNamespace(index int) (ns string, err error)

	// -- GetAttributeName ------------------------------------------

	// Returns the local name of the specified attribute
	// if namespaces are enabled or just attribute name if namespaces are disabled.
	// Throws an IndexOutOfBoundsException if the index is out of range
	// or current event type is not START_TAG.
	//
	// @param zero based index of attribute
	// @return attribute names
	//
	GetAttributeName(index int) string

	// -- GetAttributePrefix ----------------------------------------

	// Returns the prefix of the specified attribute
	// Returns null if the element has no prefix.
	// If namespaces are disabled it will always return null.
	// Returns IndexOutOfBounds if the index is out of range
	// or current event type is not START_TAG.
	//
	// @param zero based index of attribute
	// @return attribute prefix or null if namesapces processing is not enabled.
	//
	GetAttributePrefix(index int) (string, error)

	// -- GetAttributeValue -----------------------------------------

	// Returns the given attribute's value.  The index is zero-based.
	// Returns an IndexOutOfBounds error if the index is out of range
	// or current event type is not START_TAG.
	//
	GetAttributeValue(index int) (string, error)

	// Returns the attributes value identified by the namespace URI and
	// namespace localName. If namespaces are disabled namespace must be empty.
	// If the current event type is not START_TAG then IndexOutOfBounds
	// error will be returned.
	//
	// "NS" added to disambiguate.
	//
	GetAttributeValueNS(namespace, name string) (string, error)

	// ACTUAL PARSING METHODS =======================================

	// -- GetEventType ----------------------------------------------

	// Returns the type of the current event (START_TAG, END_TAG, TEXT, etc.)
	//
	// @see #next()
	// @see #nextToken()
	// May throw XmlPullParserException
	//
	GetEventType() (PullEvent, error)

	// -- Next ------------------------------------------------------

	// Get next parsing event - element content wil be coalesced and only one
	// TEXT event must be returned for whole element content (comments and
	// processing instructions will be ignored and emtity references
	// must be expanded or exception mus be thrown if entity reerence cannot
	// be exapnded). If element content is empty (content is "") then no
	// TEXT event will be reported.
	//
	// <p><b>NOTE:</b> empty element (such as &lt;tag/>) will be reported
	//  with  two separate events: START_TAG, END_TAG - it must be so to preserve
	//   parsing equivalency of empty element to &lt;tag>&lt;/tag>.
	//  (see isEmptyElementTag ())
	//
	// @see #isEmptyElementTag
	// @see #START_TAG
	// @see #TEXT
	// @see #END_TAG
	// @see #END_DOCUMENT

	// May throw XmlPullParserException or IOError
	Next() (PullEvent, error)

	// -- NextToken -------------------------------------------------

	// This method works similarly to next() but will expose
	// additional event types (COMMENT, CDSECT, DOCDECL, ENTITY_REF,
	// PROCESSING_INSTRUCTION, or IGNORABLE_WHITESPACE) if they are
	// available in input.
	//
	// <p>If special feature FEATURE_XML_ROUNDTRIP
	// (identified by URI: http://xmlpull.org/v1/doc/features.html#xml-roundtrip)
	// is true then it is possible to do XML document round trip ie. reproduce
	// exectly on output the XML input using getText().
	//
	// <p>Here is the list of tokens that can be  returned from nextToken()
	// and what getText() and getTextCharacters() returns:<dl>
	// <dt>START_DOCUMENT<dd>null
	// <dt>END_DOCUMENT<dd>null
	// <dt>START_TAG<dd>null
	//   unless FEATURE_XML_ROUNDTRIP enabled and then returns XML tag, ex: &lt;tag attr='val'>
	// <dt>END_TAG<dd>null
	// unless FEATURE_XML_ROUNDTRIP enabled and then returns XML tag, ex: &lt;/tag>
	// <dt>TEXT<dd>return unnormalized element content
	// <dt>IGNORABLE_WHITESPACE<dd>return unnormalized characters
	// <dt>CDSECT<dd>return unnormalized text <em>inside</em> CDATA
	//  ex. 'fo&lt;o' from &lt;!CDATA[fo&lt;o]]>
	// <dt>PROCESSING_INSTRUCTION<dd>return unnormalized PI content ex: 'pi foo' from &lt;?pi foo?>
	// <dt>COMMENT<dd>return comment content ex. 'foo bar' from &lt;!--foo bar-->
	// <dt>ENTITY_REF<dd>return unnormalized text of entity_name (&entity_name;)
	// <br><b>NOTE:</b> it is user responsibility to resolve entity reference
	// <br><b>NOTE:</b> character entities and standard entities such as
	//  &amp;amp; &amp;lt; &amp;gt; &amp;quot; &amp;apos; are reported as well
	// and are not resolved and not reported as TEXT tokens!
	// This requirement is added to allow to do roundtrip of XML documents!
	// <dt>DOCDECL<dd>return inside part of DOCDECL ex. returns:<pre>
	// &quot; titlepage SYSTEM "http://www.foo.bar/dtds/typo.dtd"
	// [&lt;!ENTITY % active.links "INCLUDE">]&quot;</pre>
	// <p>for input document that contained:<pre>
	// &lt;!DOCTYPE titlepage SYSTEM "http://www.foo.bar/dtds/typo.dtd"
	// [&lt;!ENTITY % active.links "INCLUDE">]></pre>
	// </dd>
	// </dl>
	//
	// NOTE: returned text of token is not end-of-line normalized.
	//
	// @see #next
	// @see #START_TAG
	// @see #TEXT
	// @see #END_TAG
	// @see #END_DOCUMENT
	// @see #COMMENT
	// @see #DOCDECL
	// @see #PROCESSING_INSTRUCTION
	// @see #ENTITY_REF
	// @see #IGNORABLE_WHITESPACE
	//
	// May throw XmlPullParserException or IOError
	NextToken() (PullEvent, error)

	// UTILITY METHODS //////////////////////////////////////////////

	// -- Require ---------------------------------------------------

	// Test if the current event is of the given type and if the
	// namespace and name do match. "" will match any namespace
	// and any name. If the current event is TEXT with isWhitespace()=
	// true, and the required type is not TEXT, next () is called prior
	// to the test. If the test is not passed, an exception is
	// thrown. The exception text indicates the parser position,
	// the expected event and the current event (not meeting the
	// requirement.
	//
	// <p>essentially it does this
	// <pre>
	//  if (getEventType() == TEXT && type != TEXT && isWhitespace ())
	//    next ();
	//
	//  if (type != getEventType()
	//  || (namespace != null && !namespace.equals (getNamespace ()))
	//  || (name != null && !name.equals (getName ()))
	//     throw new XmlPullParserException ( "expected "+ TYPES[ type ]+getPositionDesctiption());
	// </pre>
	//
	// May throw XmlPullParserException or IOError
	//
	Require(type_ PullEvent, namespace, name string)

	// -- ReadText --------------------------------------------------

	// If the current event is text, the value of getText is
	// returned and next() is called. Otherwise, an empty
	// String ("") is returned. Useful for reading element
	// content without needing to performing an additional
	// check if the element is empty.
	//
	// <p>essentially it does this
	// <pre>
	//   if (getEventType != TEXT) return ""
	//   String result = getText ();
	//   next ();
	//   return result;
	// </pre>
	//
	// May throw XmlPullParserException or IOError
	//
	ReadText() string
}
