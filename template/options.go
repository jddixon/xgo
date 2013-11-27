package template

// xgo/template/options.go

import (
	gc "github.com/jddixon/xgo/context"
)

// Options normally set from the command line or derived from those.
// Not used in this package but used by xlReg
type Options struct {
	BDir      string
	Context   *gc.Context
	InputExt  string
	FileNames []string
	JustShow  bool
	OutputExt string
	TDir      string
	Testing   bool
	Verbose   bool
}
