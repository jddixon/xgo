package om

// xgo/xml/om/cycle_checker.go

import ()

// Walk down the subtrees, child nodes and their descendents,
// looking for this node; if found, there is a cycle in the graph.
//
// This implements VisitorI
//
type CycleChecker struct {
	ThisHolder *Holder
}

func NewCycleChecker(h *Holder) (cc *CycleChecker, err error) {
	if h == nil {
		err = NilHolder
	} else {
		cc = CycleChecker{
			ThisHolder: h,
		}
	}
	return
}

// On arriving at the node, do the identity check.//
func (cc *CycleChecker) OnEntry(node Node) (err error) {
	if node == cc.ThisHolder {
		return GraphCycleError
	}
}

// On leaving, do nothing//
func (cc *CycleChecker) OnExit(Node node) (err error) {
	return
}
