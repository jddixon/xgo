package template

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _ = fmt.Print

// Given a source directory (tDir), a list of base file names, the input
// extension inputExt, a target directory(bDir), and a target extension
// outputExt, and a context, process each file, replacing each ${SYMBOL}
// with the value specified in the context.
//
// The caller guarantees that tDir exists and that the list of file names
// is non-empty.
func Process(options *Options) (err error) {
	var (
		in, out *os.File
	)
	context := options.Context
	tDir := options.TDir // path to source/template directory
	bDir := options.BDir // path to output/build directory
	fileNames := options.FileNames
	inputExt := options.InputExt
	outputExt := options.OutputExt
	verbose := options.Verbose

	// if bDir does not exist, create it
	err = os.MkdirAll(bDir, 0755)

	for i := 0; i < len(fileNames); i++ {
		inName := fileNames[i] + inputExt
		pathToIn := filepath.Join(tDir, inName)
		outName := fileNames[i] + outputExt
		pathToOut := filepath.Join(bDir, outName)
		in, err = os.Open(pathToIn)
		if err == nil {
			out, err = os.Create(pathToOut)
			if err == nil {
				t, err := NewTemplate(in, out, context)
				if err == nil {
					err = t.Apply()
				}
			}
		}
	}

	_ = verbose
	return
}

// Now superfluous.
func ProcessFile(options *Options, in io.Reader, out io.Writer) (err error) {

	t, err := NewTemplate(in, out, options.Context)
	if err == nil {
		err = t.Apply()
	}
	return
}
