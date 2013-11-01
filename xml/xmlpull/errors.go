package xmlpull

import (
	e "errors"
)

// The last is of course a Java name, and will be changed or will be replaced
// by a more specific error such as UnsupportedFeature.
var (
	EmptyFeatureString       = e.New("XmlPullParser empty feature string")
	IllegalArgument          = e.New("XmlPullParser illegal argument")
	IndexOutOfBounds         = e.New("XmlPullParser index out of bounds")
	IOError                  = e.New("XmlPullParser io error")
	NotOnStartTag            = e.New("XmlPullParser not on start tag")
	NotElementContentRelated = e.New("XmlPullParser not element-content-related")
	PosOutOfRange            = e.New("XmlPullParser pos out of range")
	UnsupportedFeature       = e.New("XmlPullParser unsupported feature")
	XmlPullParserException   = e.New("XmlPullParser exception")
)
