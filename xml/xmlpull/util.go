package xmlpull

// Compare two slices.  If either is nil, return false.  If the lengths
// differ or any of the constituent runes differ, return false.
// Otherwise return true.
func SameRunes(a, b []rune) (same bool) {
	same = true
	if a == nil || b == nil {
		same = false
	} else {
		lenA := len(a)
		lenB := len(b)
		if lenA != lenB {
			same = false
		} else {
			for i := 0; i < lenA; i++ {
				if a[i] != b[i] {
					same = false
					break
				}
			}
		}
	}
	return
}
