package xmlpull

// Simplistic implementation of hash function that has constant time to
// compute - so it also means diminishing hash quality for long strings
// but for XML parsing it should be good enough ...

func FastHash(runes []rune) (hash uint32) {

	length := len(runes)
	if length > 0 {
		rHash := runes[0] // first char is rHash
		rHash = (rHash << 7) + runes[length-1]
		if length > 16 {
			rHash = (rHash << 7) + runes[length/4] // 1/4 from beginning
		}
		if length > 8 {
			rHash = (rHash << 7) + runes[length/2] // 1/2 of string size ...
		}
		// Hash is at most done 3 times <<7 so shifted by 21 bits 8 bit value
		// so max result == 29 bits so it is quite just below 31 bits for uint32
		hash = uint32(rHash)
	}
	return hash
}

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
