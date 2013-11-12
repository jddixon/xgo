package om

// xgo/xml/om/nodeI.go

type NodeI interface {
	GetDocument() *Document
	SetDocument(newDoc *Document) error
	GetHolder() *Element
	SetHolder(h *Element)
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
