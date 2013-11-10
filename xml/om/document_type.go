package om

// xgo/xml/om/document_type.go

import (
	"fmt"
)

type DocumentType struct {
	name, value string
	Node
}

// This unfortunately makes no sense as-is.
func NewDocumentType(name, value string) (dt *DocumentType) {
	// super();
	dt = &DocumentType{
		name:  name,
		value: value,
	}
	return
}

// PROPERTIES ///////////////////////////////////////////////////
func (dt *DocumentType) GetName() string {
	return dt.name
}
func (dt *DocumentType) GetValue() string {
	return dt.value
}

// NODE METHODS /////////////////////////////////////////////////
func (dt *DocumentType) IsDocType() bool {
	return true
}

func (dt *DocumentType) ToXml() string {
	return fmt.Sprintf(" %s=\"%s\"", dt.name, dt.value)
}
