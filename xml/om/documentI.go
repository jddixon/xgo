package om

// xgo/xml/om/documentI.go

type DocumentI interface {
	GetDocType() *DocumentType
	SetDocType(dt *DocumentType) error
	GetEncoding() string
	GetElementNode() ElementI
	SetElementNode(elm ElementI) error
	GetVersion() string
	ElementI
}
