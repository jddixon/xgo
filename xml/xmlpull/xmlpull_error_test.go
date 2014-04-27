package xmlpull

// xgo/xml/xmlpull/xmlpull_error_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
	"io"
	"strings"
)

var _ = fmt.Print

func (s *XLSuite) TestXmlPullError(c *C) {

	const (
		ERR_MSG_1 = "test msg 1"
		ERR_MSG_2 = "second test msg"
		ERR_MSG_3 = "error msg 3"
	)
	var rd1 io.Reader = strings.NewReader("abc\ndef\nversion 97.1 ")
	xpp, err := NewNewParser(rd1) // accept default encoding
	c.Assert(err, IsNil)
	c.Assert(xpp, NotNil)

	err = xpp.ExpectStr("abc")
	c.Assert(err, IsNil)

	c.Assert(xpp.LineNo(), Equals, 1)
	c.Assert(xpp.ColNo(), Equals, 3)

	e1 := xpp.NewXmlPullError(ERR_MSG_1)
	expected := fmt.Sprintf("line %d col %d: %s",
		xpp.LineNo(), xpp.ColNo(), ERR_MSG_1)
	c.Assert(e1.Error(), Equals, expected)

	nl, err := xpp.NextCh()
	c.Assert(err, IsNil)
	c.Assert(nl, Equals, '\n')
	c.Assert(xpp.LineNo(), Equals, 2)
	c.Assert(xpp.ColNo(), Equals, 0)

	e2 := xpp.NewXmlPullError(ERR_MSG_2)
	expected = fmt.Sprintf("line %d col %d: %s",
		xpp.LineNo(), xpp.ColNo(), ERR_MSG_2)
	c.Assert(e2.Error(), Equals, expected)

	err = xpp.ExpectStr("def")
	c.Assert(err, IsNil)
	c.Assert(xpp.LineNo(), Equals, 2)
	c.Assert(xpp.ColNo(), Equals, 3)

	err = xpp.ExpectCh('\n')
	c.Assert(err, IsNil)
	c.Assert(xpp.LineNo(), Equals, 3)
	c.Assert(xpp.ColNo(), Equals, 0)

	err = xpp.ExpectStr("version 97.1 ")
	c.Assert(err, IsNil)
	c.Assert(xpp.LineNo(), Equals, 3)
	c.Assert(xpp.ColNo(), Equals, 13)

	e3 := xpp.NewXmlPullError(ERR_MSG_3)
	expected = fmt.Sprintf("line %d col %d: %s",
		xpp.LineNo(), xpp.ColNo(), ERR_MSG_3)
	c.Assert(e3.Error(), Equals, expected)

}
