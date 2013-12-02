package om

// xgo/xml/om/nodeI.go

type NodeI interface {
	GetDocument() DocumentI
	GetHolder() ElementI
	SetDocument(newDoc DocumentI) error
	SetHolder(h ElementI)
	ToXml() string
	WalkAll(VisitorI) error

	// value is questsionable
	IsAttr() bool
	IsCdata() bool
	IsComment() bool
	IsDocType() bool
	IsElement() bool
	IsProcessingInstruction() bool
	IsText() bool
}
