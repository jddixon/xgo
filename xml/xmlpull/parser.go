package xmlpull

import(
	gu "github.com/jddixon/xgo/util"
)

type Parser struct {
	si	gu.StrIntern


}

func NewParser() (p *Parser, err error) {

	si := gu.NewStrIntern()


	if err == nil {
		p = &Parser {
			si:	si,
		}
	}
	return
}


