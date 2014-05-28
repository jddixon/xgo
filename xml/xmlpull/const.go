package xmlpull

// xgo/xml/xmlpull/const.go

var (
	// The default namespace.
	NO_NAMESPACE = "" //  make([]rune, 0)
)

const (

	// Namespace-related features

	// FEATURE: Processing of namespaces is by default set to false.
	//
	// Cannot be changed during parsing.
	FEATURE_PROCESS_NAMESPACES = "http://xmlpull.org/v1/doc/features.html#process-namespaces"

	// FEATURE: Report namespace attributes also - they can be distinguished
	// looking for prefix == "xmlns" or prefix == null and name == "xmlns.
	// Off by default and only meaningful when FEATURE_PROCESS_NAMESPACES
	// feature is on.
	//
	// Cannot be changed during parsing.
	FEATURE_REPORT_NAMESPACE_ATTRIBUTES = "http://xmlpull.org/v1/doc/features.html#report-namespace-prefixes"

	// FEATURE: Processing of DOCDECL is by default set to false
	// and if DOCDECL is encountered it is reported by NextToken()
	// and ignored by Next().
	//
	// If processing is set to true then DOCDECL must be processed by parser.
	//
	// NOTE: if the DOCDECL was ignored further in parsing there may be fatal
	// exception (panic) when undeclared entity is encountered.
	//
	// Cannot be changed during parsing.
	FEATURE_PROCESS_DOCDECL = "http://xmlpull.org/v1/doc/features.html#process-docdecl"

	// FEATURE: Report all validation errors as defined by XML 1.0
	// specification (implies that FEATURE_PROCESS_DOCDECL is true and both
	// internal and external DOCDECL will be processed).
	//
	// Cannot be changed during parsing.
	FEATURE_VALIDATION = "http://xmlpull.org/v1/doc/features.html#validation"
)

type PullEvent int

const (
	// ==============================================================
	// EVENT TYPES AS REPORTED BY Next()
	// ==============================================================

	// EVENT TYPE and TOKEN: signalize that parser is at the very beginning
	// of the document and nothing has been read yet - the parser is before
	// the first call to Next() or NextToken() (available from
	// <a href="#Next()">Next()</a> and <a href="#NextToken()">NextToken()</a>).
	//
	// @see #Next
	// @see #NextToken
	//
	START_DOCUMENT PullEvent = iota

	// EVENT TYPE and TOKEN: logical end of xml document
	// (available from
	// <a href="#Next()">Next()</a> and <a href="#NextToken()">NextToken()</a>).
	//
	// <p><strong>NOTE:</strong> calling again
	// <a href="#Next()">Next()</a> or <a href="#NextToken()">NextToken()</a>
	// will result in exception being thrown.
	//
	// @see #Next
	// @see #NextToken
	//
	END_DOCUMENT

	// EVENT TYPE and TOKEN: start tag was just read
	// (available from <a href="#Next()">Next()</a> and
	// <a href="#NextToken()">NextToken()</a>).
	// The name of start tag is available from GetName(), its namespace and
	// prefix are available from GetNamespace() and GetPrefix()
	// if <a href='#FEATURE_PROCESS_NAMESPACES'>namespaces are enabled</a>.
	// See getAttribute* methods to retrieve element attributes.
	// See GetNamespace* methods to retrieve newly declared namespaces.
	//
	// @see #Next
	// @see #NextToken
	// @see #GetName
	// @see #GetPrefix
	// @see #GetNamespace
	// @see #getAttributeCount
	// @see #getDepth
	// @see #GetNamespaceCount
	// @see #GetNamespace
	// @see #FEATURE_PROCESS_NAMESPACES
	//
	START_TAG

	// EVENT TYPE and TOKEN: end tag was just read (available from
	// <a href="#Next()">Next()</a> and <a href="#NextToken()">NextToken()</a>).
	// The name of start tag is available from GetName(), its namespace and
	// prefix are available from GetNamespace() and GetPrefix()
	//
	// @see #Next
	// @see #NextToken
	// @see #GetName
	// @see #GetPrefix
	// @see #GetNamespace
	// @see #FEATURE_PROCESS_NAMESPACES
	//
	END_TAG

	// EVENT TYPE and TOKEN: character data was read and will be available by
	// a call to GetText() (available from <a href="#Next()">Next()</a> and
	// <a href="#NextToken()">NextToken()</a>).
	// <p><strong>NOTE:</strong> Next() will (in contrast to NextToken ())
	// accumulate multiple events into one TEXT event, skipping
	// IGNORABLE_WHITESPACE, PROCESSING_INSTRUCTION and COMMENT events.
	// <p><strong>NOTE:</strong> if state was reached by calling Next() the
	// text value will be normalized and if the token was returned by
	// NextToken() then GetText() will return unnormalized content (no
	// end-of-line normalization - it is content exactly as in input XML)
	//
	// @see #Next
	// @see #NextToken
	// @see #GetText
	//
	TEXT

	// ==============================================================
	// ADDITIONAL EVENTS EXPOSED BY LOWER-LEVEL NextToken()
	// ==============================================================

	// TOKEN: CDATA sections was just read (this token is available only
	// from <a href="#NextToken()">NextToken()</a>).  The value of text
	// inside CDATA section is available  by callling GetText().
	//
	// @see #NextToken
	// @see #GetText
	//
	CDSECT

	// TOKEN: Entity reference was just read (this token is available only
	// from <a href="#NextToken()">NextToken()</a>).
	// The entity name is available by calling GetText() and it is the user's
	// responsibility to resolve entity references.
	//
	// @see #NextToken
	// @see #GetText
	//
	ENTITY_REF

	// TOKEN: Ignorable whitespace was just read (this token is available
	// only from <a href="#NextToken()">NextToken()</a>).
	// For non-validating parsers, this event is only reported by NextToken()
	// when outside the root elment.
	// Validating parsers may be able to detect ignorable whitespace at
	// other locations.
	// The value of ignorable whitespace is available by calling GetText()
	//
	// <p><strong>NOTE:</strong> this is different from calling the
	// IsWhitespace() method as element content may be whitespace but may
	// not be ignorable whitespace.
	//
	// @see #NextToken
	// @see #GetText
	//
	IGNORABLE_WHITESPACE

	// TOKEN: XML processing instruction declaration was just read
	// and GetText() will return text that is inside processing instruction
	// (this token is available only from
	// <a href="#NextToken()">NextToken()</a>).
	//
	// @see #NextToken
	// @see #GetText
	//
	PROCESSING_INSTRUCTION

	// TOKEN: XML comment was just read and GetText() will return value
	// inside the comment (this token is available only from
	// <a href="#NextToken()">NextToken()</a>).
	//
	// @see #NextToken
	// @see #GetText
	//
	COMMENT

	// TOKEN: XML DOCTYPE declaration was just read and GetText() will
	// return any text inside the DOCDECL (this token is available only
	// from <a href="#NextToken()">NextToken()</a>).
	//
	// @see #NextToken
	// @see #GetText
	//
	DOCDECL
)

// Use this array to convert evebt type number (such as START_TAG) to
// to string giving event name, ex: "START_TAG" == TYPES[START_TAG]
//
var PULL_EVENT_NAMES = []string{
	"START_DOCUMENT",
	"END_DOCUMENT",
	"START_TAG",
	"END_TAG",
	"TEXT",
	"CDSECT",
	"ENTITY_REF",
	"IGNORABLE_WHITESPACE",
	"PROCESSING_INSTRUCTION",
	"COMMENT",
	"DOCDECL",
}

// parser states as used in Next() ----------------------------------
// XXX STATES AND NAMES MUST BE KEPT SYNCHRONIZED

type ParserState uint

const (
	PRE_START_DOC ParserState = iota
	START_STATE
	XML_DECL_SEEN
	DOC_DECL_SEEN
	START_ROOT_SEEN

	END_ROOT_SEEN
	PAST_END_DOC
)

var PARSER_STATE_NAMES = []string{
	"PRE_START_DOC",
	"START_STATE",
	"XML_DECL_SEEN",
	"DOC_DECL_SEEN",
	"START_ROOT_SEEN",

	"END_ROOT_SEEN",
	"PAST_END_DOC",
}
