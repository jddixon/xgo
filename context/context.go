package context

import (
	"fmt"
	"strings"
	"sync"
)

// A naming context consisting of a possibly nested set of name-to-object
// bindings.  If there is a parent context and a key cannot be resolved
// in this context, an attempt will be made to resolve it in the parent,
// recursively.
//
// Names added to the context must not be null.
//
// This implementation is intended to be thread-safe.
//
type Context struct {
	ctx    map[string]interface{}
	parent *Context
	mu     sync.RWMutex
}

// CONSTRUCTORS /////////////////////////////////////////////////

// Create a context without a Parent.
//
func NewNewContext() *Context {
	return &Context{
		ctx: make(map[string]interface{}),
	}
}

// Create a context with a Parent Context.  The parent may be nil.
//
func NewContext(parent *Context) (k *Context) {
	k = NewNewContext()
	k.parent = parent
	return
}

// METHODS //////////////////////////////////////////////////////

// Bind a name to an object at this Context level.  Neither name
// nor object may be null.
//
// If this context has a Parent, the binding at this level will
// mask any bindings in the Parent and above.
//
func (k *Context) Bind(name string, value interface{}) (err error) {
	if name == "" {
		err = EmptyName
	} else if value == nil {
		err = NilValue
	} else {
		k.mu.Lock()
		k.ctx[name] = value
		k.mu.Unlock()
	}
	return
}

// Looks up a name recursively.  If the name is bound at this level,
// the object it is bound to is returned.  Otherwise, if there is
// a Parent Context, the value returned by a lookup in the Parent
// Context is returned.  If there is no Parent and no match, returns
// nil.
//
func (k *Context) Lookup(name string) (value interface{}, err error) {
	if name == "" {
		err = EmptyName
	} else {
		k.mu.RLock()
		defer k.mu.RUnlock()
		value = k.ctx[name]
		if value == nil && k.parent != nil {
			value, err = k.parent.Lookup(name)
		}
	}
	return
}

// Remove a binding from the context.  If there is no such binding,
// silently ignore the request.  Any binding at a higher level, in
// the Parent Context or above, is unaffected by this operation.
//
func (k *Context) Unbind(name string) (err error) {
	if name == "" {
		err = EmptyName
	} else {
		k.mu.Lock()
		delete(k.ctx, name)
		k.mu.Unlock()
	}
	return
}

// PROPERTIES ///////////////////////////////////////////////////

// Return the number of bindings at this level.
//
func (k *Context) Size() int {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return len(k.ctx)
}

// Return a reference to the Parent Context or nil if there is none
//
func (k *Context) GetParent() *Context {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.parent
}

// Change the Parent Context. This method returns a reference to
// this instance, to allow method calls to be chained.
//
func (k *Context) SetParent(newParent *Context) *Context {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.parent = newParent
	return k
}

// Serialize a map[string]string context.
//
// XXX This will not handle a map[string]interface{}.
// XXX Missing data integrity hash.
//
func (k *Context) String() string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	var ss []string
	for k, v := range k.ctx {
		s := fmt.Sprintf("%s\t%s\n", k, v)
		ss = append(ss, s)
	}
	return strings.Join(ss, "")
}

// Parse a context which has been serialized using Context.String().
//
// XXX This will not handle a map[string]interface{}.
// XXX Does not handle any data integrity hash.
//
func ParseContext(s string) (k *Context, err error) {
	if s == "" {
		err = EmptySerialization
	} else {
		ss := strings.Split(s, "\n")
		m := make(map[string]interface{})
		for i := 0; i < len(ss); i++ {
			s := strings.TrimSpace(ss[i])
			if len(s) == 0 || strings.HasPrefix(s, "#") {
				continue
			}
			parts := strings.Split(s, "\t")
			if len(parts) != 2 {
				err = IllFormedCtxSerialization
				break
			} else {
				k := strings.TrimRight(parts[0], " ")
				v := strings.TrimLeft(parts[1], " ")
				m[k] = v
			}
		}
		if err == nil {
			k = &Context{ctx: m}
		}
	}
	return
}
