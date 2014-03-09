package xmlpull

import (
	"fmt"
	"strings"
)

func printableChar(ch rune) (out string) {
	if ch == '\n' {
		out = "\\n"
	} else if ch == '\r' {
		out = "\\r"
	} else if ch == '\t' {
		out = "\\t"
	} else if ch == '\'' {
		out = "\\'"
	} else if ch < 32 || 127 < ch {
		out = fmt.Sprintf("\\u%x", ch)
	} else {
		out = string(ch) // GEEP
	}
	return
}

func printableStr(s string) (out string) {
	if s != "" {
		runes := []rune(s)
		var ss []string
		for i := 0; i < len(s); i++ {
			ss = append(ss, printableChar(runes[i]))
		}
		out = strings.Join(ss, "")
	}
	return
}
