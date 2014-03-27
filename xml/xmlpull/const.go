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
	// and if DOCDECL is encountered it is reported by nextToken()
	// and ignored by next().
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

	// EVENT TYPE and TOKEN: signalize that parser is at the very beginning of the document
	// and nothing was read yet - the parser is before first call to next() or nextToken()
	// (available from <a href="#next()">next()</a> and <a href="#nextToken()">nextToken()</a>).
	//
	// @see #next
	// @see #nextToken
	//
	START_DOCUMENT PullEvent = iota

	// EVENT TYPE and TOKEN: logical end of xml document
	// (available from <a href="#next()">next()</a> and <a href="#nextToken()">nextToken()</a>).
	//
	// <p><strong>NOTE:</strong> calling again
	// <a href="#next()">next()</a> or <a href="#nextToken()">nextToken()</a>
	// will result in exception being thrown.
	//
	// @see #next
	// @see #nextToken
	//
	END_DOCUMENT

	// EVENT TYPE and TOKEN: start tag was just read
	// (available from <a href="#next()">next()</a> and <a href="#nextToken()">nextToken()</a>).
	// The name of start tag is available from getName(), its namespace and prefix are
	// available from getNamespace() and getPrefix()
	// if <a href='#FEATURE_PROCESS_NAMESPACES'>namespaces are enabled</a>.
	// See getAttribute* methods to retrieve element attributes.
	// See getNamespace* methods to retrieve newly declared namespaces.
	//
	// @see #next
	// @see #nextToken
	// @see #getName
	// @see #getPrefix
	// @see #getNamespace
	// @see #getAttributeCount
	// @see #getDepth
	// @see #getNamespaceCount
	// @see #getNamespace
	// @see #FEATURE_PROCESS_NAMESPACES
	//
	START_TAG

	// EVENT TYPE and TOKEN: end tag was just read
	// (available from <a href="#next()">next()</a> and <a href="#nextToken()">nextToken()</a>).
	// The name of start tag is available from getName(), its namespace and prefix are
	// available from getNamespace() and getPrefix()
	//
	// @see #next
	// @see #nextToken
	// @see #getName
	// @see #getPrefix
	// @see #getNamespace
	// @see #FEATURE_PROCESS_NAMESPACES
	//
	END_TAG

	// EVENT TYPE and TOKEN: character data was read and will be available by call to getText()
	// (available from <a href="#next()">next()</a> and <a href="#nextToken()">nextToken()</a>).
	// <p><strong>NOTE:</strong> next() will (in contrast to nextToken ()) accumulate multiple
	// events into one TEXT event, skipping IGNORABLE_WHITESPACE,
	// PROCESSING_INSTRUCTION and COMMENT events.
	// <p><strong>NOTE:</strong> if state was reached by calling next() the text value will
	// be normalized and if the token was returned by nextToken() then getText() will
	// return unnormalized content (no end-of-line normalization - it is content exactly as in
	// input XML)
	//
	// @see #next
	// @see #nextToken
	// @see #getText
	//
	TEXT

	// ==============================================================
	// ADDITIONAL EVENTS EXPOSED BY LOWER-LEVEL nextToken()
	// ==============================================================

	// TOKEN: CDATA sections was just read
	// (this token is available only from <a href="#nextToken()">nextToken()</a>).
	// The value of text inside CDATA section is available  by callling getText().
	//
	// @see #nextToken
	// @see #getText
	//
	CDSECT

	// TOKEN: Entity reference was just read
	// (this token is available only from <a href="#nextToken()">nextToken()</a>).
	// The entity name is available by calling getText() and it is user responsibility
	// to resolve entity reference.
	//
	// @see #nextToken
	// @see #getText
	//
	ENTITY_REF

	// TOKEN: Ignorable whitespace was just read
	// (this token is available only from <a href="#nextToken()">nextToken()</a>).
	// For non-validating
	// parsers, this event is only reported by nextToken() when
	// outside the root elment.
	// Validating parsers may be able to detect ignorable whitespace at
	// other locations.
	// The value of ignorable whitespace is available by calling getText()
	//
	// <p><strong>NOTE:</strong> this is different than callinf isWhitespace() method
	//    as element content may be whitespace but may not be ignorable whitespace.
	//
	// @see #nextToken
	// @see #getText
	//
	IGNORABLE_WHITESPACE

	// TOKEN: XML processing instruction declaration was just read
	// and getText() will return text that is inside processing instruction
	// (this token is available only from <a href="#nextToken()">nextToken()</a>).
	//
	// @see #nextToken
	// @see #getText
	//
	PROCESSING_INSTRUCTION

	// TOKEN: XML comment was just read and getText() will return value inside comment
	// (this token is available only from <a href="#nextToken()">nextToken()</a>).
	//
	// @see #nextToken
	// @see #getText
	//
	COMMENT

	// TOKEN: XML DOCTYPE declaration was just read
	// and getText() will return text that is inside DOCDECL
	// (this token is available only from <a href="#nextToken()">nextToken()</a>).
	//
	// @see #nextToken
	// @see #getText
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
