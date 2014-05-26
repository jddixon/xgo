package xmlpull

import (
	gl "github.com/jddixon/xgo/lex"
	"io"
)

const (
	READ_CHUNK_SIZE = 4096 // or multiple
)

// Any fields added here also should be added to reset()
//
type Parser struct {
	// VETTED FIELDS ////////////////////////////////////////////////
	curEvent PullEvent //  defined in const.go

	// defined in the prolog
	xmlDeclVersion, xmlDeclEncoding string
	xmlDeclContent                  string
	xmlDeclStandalone               bool
	docTypeDecl                     string

	// doNext State -------------------------------------------------
	isEmptyElement bool
	pastEndTag     bool
	rootElmSeen    bool
	seenStartTag   bool

	// OTHER FIELDS /////////////////////////////////////////////////

	// parser behavior
	normalizeIgnorableWS bool
	processNamespaces    bool
	roundtripSupported   bool

	startLine, startCol int // where a syntactic element begins

	// accumulated characters of various types -- yes, kludgey
	cDataChars   []rune
	commentChars []rune
	piChars      []rune
	piTarget     []rune

	// transient variables for NextToken()
	tokenizing    bool
	text          []rune
	entityRefName string // []rune

	// global parser state
	lineNo int // line number		// redundant
	colNo  int // column number		// redundant

	location      string
	reachedEnd    bool // used only in parseEpilog?
	seenEndTag    bool
	seenAmpersand bool
	afterLT       bool // have encountered a left angle bracket (<)
	seenDocdecl   bool

	// element stack
	elmDepth      uint
	elRawName     [][]rune
	elRawNameLine []int

	elName           []string // []rune
	elPrefix         []string // []rune
	elUri            []string // [][]rune
	elValue          [][]rune
	elNamespaceCount []uint

	// attribute stack
	attributeCount    uint
	attributeName     []string //[]rune
	attributeNameHash []uint32 // stores FastHash output
	attributePrefix   []string // []rune
	attributeUri      []string // [][]rune
	attributeValue    []string // []rune

	// namespace stack
	namespaceCount      uint
	namespacePrefix     []string //[]rune
	namespacePrefixHash []uint32 // more FastHash output
	namespaceUri        []string // []rune

	namespaceEnd uint

	// entity replacement stack ---------------------------
	entityEnd            int
	entityName           []string // [][]rune
	entityNameBuf        [][]rune
	entityReplacement    []string // ][]rune
	entityReplacementBuf [][]rune
	entityNameHash       []uint32

	// XXX hmmm
	inputEncoding string
	gl.LexInput
}

// Called at the beginning of a syntactic construct.
func (p *Parser) start() {
	p.startLine = p.LineNo()
	p.startCol = p.ColNo()
}

// Return an XmlPullParser with the default encoding
func NewNewParser(reader io.Reader) (*Parser, error) {
	return NewParser(reader, "")
}

func NewParser(reader io.Reader, encoding string) (p *Parser, err error) {

	var lx *gl.LexInput
	if reader == nil {
		err = NilReader
	} else {
		lx, err = gl.NewLexInput(reader, encoding)
	}

	if err == nil {
		p = &Parser{
			LexInput: *lx,
		}
		p.reset()
	}
	return
}

// Return a pointer to the parser's lexer.
func (xpp *Parser) GetLexer() *gl.LexInput {
	return &xpp.LexInput
}

// Reset all global variables; also used to set initial values at
// program start.
//
// All fields in the Parser struct should be reinitialized here.

func (xpp *Parser) reset() {
	// VETTED FIELDS ------------------------------------------------
	xpp.curEvent = START_DOCUMENT

	// OTHER FIELDS -------------------------------------------------
	xpp.afterLT = false
	xpp.colNo = 0
	xpp.elmDepth = 0
	xpp.lineNo = 1
	xpp.namespaceCount = 0

	// XXX STUB XXX

}
