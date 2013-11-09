package util

// xgo/util/str_intern.go

import (
	"sync"
	//"unsafe"
)

// Use this to ensure that there is only one copy of a string in memory.
//
// See the discussion on stackoverflow, 13017449,
// is-there-an-equivalent-to-javas-string-intern-function-in-go
//
type StrIntern struct {
	m  map[string]string
	mu *sync.RWMutex
}

func NewStrIntern() StrIntern {
	return StrIntern{
		m:  make(map[string]string),
		mu: &sync.RWMutex{},
	}
}

func (si StrIntern) Intern(s string) (out string) {
	var ok bool
	si.mu.RLock()
	if out, ok = si.m[s]; ok {
		// it's already interned
		si.mu.RUnlock()
		return
	}
	si.mu.RUnlock()

	// It hasn't been interned.  So get a write lock, check again, and
	// intern it unless (very unlikely!) someone else has gotten there first.
	b := []byte(s)
	si.mu.Lock()
	defer si.mu.Unlock()
	if out, ok = si.m[s]; ok {
		// Somebody just beat us to it.
		return
	} else {
		// We got the write lock and it's still not in the cache.

		// The next two lines are equivalent.  The second should be faster.
		out = string(b) // version 1
		// out = *(*string)(unsafe.Pointer(&b))		// version 2
		si.m[out] = out
	}
	return
}
