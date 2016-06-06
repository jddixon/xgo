package template

// xgo/template/apply_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	gc "github.com/jddixon/xgo_go/context"
	xf "github.com/jddixon/xlUtil_go/lfs"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = fmt.Print

func (s *XLSuite) makeSymbolSet(c *C, rng *xr.PRNG, size int) (ss []string) {
	// make size unique names
	ss = make([]string, size)
	nameCache := make(map[string]string)
	for i := 0; i < size; i++ {
		name := rng.NextFileName(8)
		_, ok := nameCache[name]
		for ok {
			name = rng.NextFileName(8)
			_, ok = nameCache[name]

		}
		nameCache[name] = ""
		ss[i] = name
	}
	return // FOO
}

func (s *XLSuite) makeContext(c *C, rng *xr.PRNG, size int) (
	k, v []string, context *gc.Context) {

	var err error
	context = gc.NewNewContext()
	k = s.makeSymbolSet(c, rng, size)
	v = make([]string, size)
	for i := 0; i < size; i++ {
		value := rng.NextFileName(8)
		v[i] = value
		err = context.Bind(k[i], value)
		c.Assert(err, IsNil)
	}
	return k, v, context
}

// Return either spaces, or a dollar sign, or a random 'word', or a
// newline, or nothing.
func (s *XLSuite) moreBits(c *C, rng *xr.PRNG) (txt string) {
	start := rng.Intn(7)
	switch start {
	case 0:
		txt += "  "
	case 1:
		txt += "$"
	case 2:
		txt += "\n"
	case 3:
		txt += rng.NextFileName(8)
	case 4:
		txt += rng.NextFileName(8)
	case 5:
		txt += rng.NextFileName(8)
	case 6: // nothing
	}
	return
}

// Create two documents, one contain quasi-random text interspersed
// with expressions like ${SYMBOL} and the other containing the same
// text but with the expressions replaced by their values.
//
func (s *XLSuite) makeDocumentPair(c *C, rng *xr.PRNG, k, v []string,
	context *gc.Context) (tDoc, bDoc string) {

	size := len(k)

	// Start with spaces, dollar, text, newline, or nothing.
	more := s.moreBits(c, rng)
	tDoc += more
	bDoc += more

	for i := 0; i < size; i++ {
		expr := "${" + k[i] + "}"
		value, err := context.Lookup(k[i])
		c.Assert(err, IsNil)
		c.Assert(value, Equals, v[i])
		tDoc += expr
		bDoc += value.(string)

		// append random content
		more = s.moreBits(c, rng)
		tDoc += more
		bDoc += more
	}
	return
}

func (s *XLSuite) TestProcessor(c *C) {

	var (
		inputExt  = ".t"
		outputExt = ".OUT"
		prefix    string // empty for now
	)
	rng := xr.MakeSimpleRNG()
	fCount := 3 + rng.Intn(5) // so 3 to 7 inclusive

	// make a scratch directory for this test run
	dirName := rng.NextFileName(8) // may already exist
	pathToDir := filepath.Join("tmp", dirName)
	found, err := xf.PathExists(pathToDir)
	c.Assert(err, IsNil)
	for found {
		dirName = rng.NextFileName(8)
		pathToDir = filepath.Join("tmp", dirName)
		found, err = xf.PathExists(pathToDir)
		c.Assert(err, IsNil)
	}
	tDir := filepath.Join(pathToDir, "t")
	bDir := filepath.Join(pathToDir, "b")

	err = os.MkdirAll(tDir, 0755)
	c.Assert(err, IsNil)
	err = os.MkdirAll(bDir, 0755)
	c.Assert(err, IsNil)

	tFiles := make([]string, fCount) // input files
	bFiles := make([]string, fCount) // output files
	baseNames := s.makeSymbolSet(c, rng, fCount)

	// create the paths to the files
	prefixLen := rng.Intn(3)
	switch prefixLen {
	case 1:
		prefix = "x"
	case 2:
		prefix = "yz"
	}
	for i := 0; i < fCount; i++ {
		fileName := baseNames[i] + inputExt
		tFiles[i] = filepath.Join(tDir, fileName)
		fileName = prefix + baseNames[i] + outputExt
		bFiles[i] = filepath.Join(bDir, fileName)
	}
	k, v, context := s.makeContext(c, rng, fCount)
	tDocs := make([]string, fCount) // input documents
	bDocs := make([]string, fCount) // output documents
	for i := 0; i < fCount; i++ {
		tDoc, bDoc := s.makeDocumentPair(c, rng, k, v, context)
		tDocs[i] = tDoc
		bDocs[i] = bDoc
		err = ioutil.WriteFile(tFiles[i], []byte(tDoc), 0666)
		c.Assert(err, IsNil)
	}

	options := &Options{
		BDir:      bDir,
		Context:   context,
		InputExt:  inputExt,
		FileNames: baseNames,
		Prefix:    prefix,
		OutputExt: outputExt,
		TDir:      tDir,
	}
	err = Process(options)
	c.Assert(err, IsNil)

	// COMPARE WHAT'S OUT THERE WITH WHAT IS EXPECTED
	for i := 0; i < fCount; i++ {
		actual, err := ioutil.ReadFile(bFiles[i])
		c.Assert(err, IsNil)
		c.Assert(string(actual), Equals, bDocs[i])
	}
}
