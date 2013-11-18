package om

// xgo/xml/om/nodeI.go

type NodeI interface {
	GetDocument() DocumentI
	GetHolder() ElementI
	SetDocument(newDoc DocumentI) error
	SetHolder(h ElementI)
	ToXml() string
	WalkAll(VisitorI) error

	// are these of any value at all?
	IsAttr() bool
	IsComment() bool
	IsDocument() bool
	IsDocType() bool
	IsElement() bool
	IsText() bool
	IsProcessingInstruction() bool
}
