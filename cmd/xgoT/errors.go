package main

import (
	e "errors"
)

var (
	NothingToDo        = e.New("nothing to do - no input files")
	SrcDirDoesNotExist = e.New("source directory does not exist")
)
