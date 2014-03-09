package xmlpull

// charTests.go

const (
	LOOKUP_MAX      = 0x400
	LOOKUP_MAX_CHAR = rune(LOOKUP_MAX)
)

var (
	lookupNameStartChar = make([]bool, LOOKUP_MAX)
	lookupNameChar      = make([]bool, LOOKUP_MAX)
)

func _setName(ch rune) {
	lookupNameChar[ch] = true
}
func _setNameStart(ch rune) {
	lookupNameStartChar[ch] = true
	_setName(ch)
}
func init() {
	_setNameStart(':')
	for ch := 'A'; ch <= 'Z'; ch++ {
		_setNameStart(ch)
	}
	_setNameStart('_')
	for ch := 'a'; ch <= 'z'; ch++ {
		_setNameStart(ch)
	}
	for ch := '\u00c0'; ch <= '\u02FF'; ch++ {
		_setNameStart(ch)
	}
	for ch := '\u0370'; ch <= '\u037d'; ch++ {
		_setNameStart(ch)
	}
	for ch := '\u037f'; ch < '\u0400'; ch++ {
		_setNameStart(ch)
	}
	_setName('-')
	_setName('.')
	for ch := '0'; ch <= '9'; ch++ {
		_setName(ch)
	}
	_setName('\u00b7')
	for ch := '\u0300'; ch <= '\u036f'; ch++ {
		_setName(ch)
	} // GEEP
}

func isNameStartChar(ch rune) bool {
	return (ch < LOOKUP_MAX_CHAR && lookupNameStartChar[ch]) ||
		(ch >= LOOKUP_MAX_CHAR && ch <= '\u2027') ||
		(ch >= '\u202A' && ch <= '\u218F') ||
		(ch >= '\u2800' && ch <= '\uFFEF')
}

func isNameChar(ch rune) bool {
	return (ch < LOOKUP_MAX_CHAR && lookupNameChar[ch]) ||
		(ch >= LOOKUP_MAX_CHAR && ch <= '\u2027') ||
		(ch >= '\u202A' && ch <= '\u218F') ||
		(ch >= '\u2800' && ch <= '\uFFEF')
}
