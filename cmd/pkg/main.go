package pkg

import (
	"fmt"
	"os"
)

const DOT_GIT = ".git"
const DOT_GIT_OBJECTS = ".git/objects"
const DOT_GIT_REFS = ".git/refs"
const DOT_GIT_HEAD = ".git/HEAD"

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", err.Error())
		os.Exit(0)
	}
}

func CheckNonIsNotExistError(err error) {
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", err.Error())
		os.Exit(0)
	}
}
