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
	GetAttributeNamespace(index int) string

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
	GetAttributePrefix(index int) string

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
	GetEventType() (int, error)

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
	Next() (int, error)

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
	NextToken() (int, error)

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
	Require(type_ int, namespace, name string)

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
