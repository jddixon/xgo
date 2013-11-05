package om

// xgo/xml/om/documentI.go

type DocumentI interface {
	GetDocType() *DocumentType
	SetDocType(dt *DocumentType) error
	GetEncoding() string
	GetElementNode() ElementI
	SetElement(elm *Element) error
	GetVersion() string
	ElementI
}
