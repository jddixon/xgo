package om

// xgo/xml/om/nodeI.go

type NodeI interface {
	GetDocument() DocumentI
	SetDocument(newDoc DocumentI) error
	GetHolder() ElementI
	SetHolder(h ElementI)
	WalkAll(VisitorI) error
	IsAttr() bool
	IsComment() bool
	IsDocument() bool
	IsDocType() bool
	IsElement() bool
	IsText() bool
	IsProcessingInstruction() bool
	ToXml() string
}
