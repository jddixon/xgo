package main

// xgo/cmd/xgoT/xgoT.go

import (
	"flag"
	"fmt"
	gt "github.com/jddixon/xgo/template"
	"os"
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}

const (
	DEFAULT_T_DIR    = "./"
	DEFAULT_B_DIR    = "./"
	DEFAULT_CTX_FILE = "context"
	DEFAULT_T_EXT    = ".t"
	DEFAULT_B_EXT    = ".go"
)

var (
	// these need to be referenced as pointers
	bDir = flag.String("b", DEFAULT_B_DIR,
		"path to build *.go files")
	ctxFile = flag.String("c", DEFAULT_CTX_FILE,
		"path to context file")
	inputExt = flag.String("e", DEFAULT_T_EXT,
		"input file extension")
	justShow = flag.Bool("j", false,
		"display option settings and exit")
	outputExt = flag.String("E", DEFAULT_B_EXT,
		"output file extension")
	tDir = flag.String("t", DEFAULT_T_DIR,
		"path to source *.t files")

	testing = flag.Bool("T", false,
		"this is a test run")
	verbose = flag.Bool("v", false,
		"be talkative")
)

func main() {
	var err error

	flag.Usage = Usage
	flag.Parse()

	// FIXUPS ///////////////////////////////////////////////////////
	// XXX STUB XXX

	// SANITY CHECKS ////////////////////////////////////////////////
	if _, err = os.Stat(*tDir); os.IsNotExist(err) {
		err = SrcDirDoesNotExist
	}
	// XXX STUB XXX

	// DISPLAY STUFF ////////////////////////////////////////////////
	if *verbose || *justShow {
		fmt.Printf("bDir         = %v\n", *bDir)
		fmt.Printf("ctxFile      = %v\n", *ctxFile)
		fmt.Printf("inputExt     = %v\n", *inputExt)
		fmt.Printf("justShow     = %v\n", *justShow)
		fmt.Printf("outputExt    = %s\n", *outputExt)
		fmt.Printf("tDir         = %s\n", *tDir)
		fmt.Printf("testing      = %v\n", *testing)
		fmt.Printf("verbose      = %v\n", *verbose)
	}

	if err != nil {
		fmt.Printf("\nerror = %s\n", err.Error())
	}
	if err != nil || *justShow {
		return
	}

	// SET UP OPTIONS ///////////////////////////////////////////////
	options := new(gt.Options)
	options.BDir = *bDir
	options.CtxFile = *ctxFile
	options.InputExt = *inputExt
	options.JustShow = *justShow
	options.OutputExt = *outputExt
	options.TDir = *tDir
	options.Testing = *testing
	options.Verbose = *verbose

	// DO USEFUL THINGS /////////////////////////////////////////////
	// XXX STUB XXX

	return
}
