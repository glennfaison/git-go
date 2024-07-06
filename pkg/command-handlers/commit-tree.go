package command_handlers

import (
	"fmt"

	"github.com/codecrafters-io/git-starter-go/pkg"
)

func CommandHandler_CommitTree(args []string) error {
	dir := "./"
	writeToFile := true

	checksum, _, err := pkg.ComputeTreeObjectForDirectory(dir, writeToFile)
	if err != nil {
		return err
	}

	fmt.Printf("%x", checksum)
	return nil
}
