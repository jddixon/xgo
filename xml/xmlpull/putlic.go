package xmlpull

// xgo/xml/xmlpull/public.go

const (
	NO_NAMESPACE = ""
    
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

type PullToken int 
const (
	START_DOCUMENT	PullToken = iota
	END_DOCUMENT 
	START_TAG 
	END_TAG 
	TEXT 
	CDSECT 
	ENTITY_REF 
	IGNORABLE_WHITESPACE 
	PROCESSING_INSTRUCTION 
	COMMENT 
	DOCDECL 
)

var TOKEN_NAMES = []string {
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
