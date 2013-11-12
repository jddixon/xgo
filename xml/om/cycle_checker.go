package om

// xgo/xml/om/cycle_checker.go

import ()

// Walk down the subtrees, child nodes and their descendents,
// looking for this node; if found, there is a cycle in the graph.
//
// This implements VisitorI
//
type CycleChecker struct {
	ThisHolder *Element
}

// The Holder-child dependency forms a directed graph.  This
// looks for cycles in that graph.
//
func NewCycleChecker(h *Element) (cc *CycleChecker, err error) {
	if h == nil {
		err = NilHolder
	} else {
		cc = &CycleChecker{
			ThisHolder: h,
		}
	}
	return
}

// On arriving at the node, do the identity check.//
func (cc *CycleChecker) OnEntry(node NodeI) (err error) {
	// XXX could be node.Equal(cc.ThisHolder)
	if node == cc.ThisHolder {
		return GraphCycleError
	}
	return
}

// On leaving, do nothing//
func (cc *CycleChecker) OnExit(node NodeI) (err error) {
	return
}
