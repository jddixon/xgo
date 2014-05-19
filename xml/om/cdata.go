package om

import ()

// An XML text node containing text which must not be subjected to
// further interpretation by the XML processor.
type Cdata struct {
	Text
}

func NewCdata(text string) (this *Cdata) {
	t := NewText(text)
	this = &Cdata{Text: *t}
	return
}

func IsCdata() bool {
	return true
}

// return the text in a CDATA wrapper
func (cd *Cdata) ToXml() string {
	return "<![CDATA[" + cd.text + "]]>\n"
}
