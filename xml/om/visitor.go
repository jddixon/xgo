package om

// Interface for classes implementing the Visitor pattern.

type VisitorI interface {

	// Action taken by a visitor on arriving at a Node.  May return
	// RuntimeError
	OnEntry(n NodeI) error

	// Action taken by the visitor on leaving the Node.  May return
	// RuntimeError
	OnExit(n NodeI) error
}
