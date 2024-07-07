package command_handlers

import (
	"flag"
	"fmt"

	"github.com/codecrafters-io/git-starter-go/pkg"
)

func WriteTree(args []string) error {
	dir := "./"
	writeToFile := true

	flag := flag.NewFlagSet("mygit write-tree", flag.ExitOnError)
	prefix := flag.String("prefix", "", "write tree object for a subdirectory <prefix>")
	flag.Parse(args)

	if prefix != nil && *prefix != "" {
		dir = *prefix
	}
	checksum, content, err := pkg.ComputeTreeObjectForDirectory(dir, writeToFile)
	if err != nil {
		return err
	}

	pkg.WriteObjectFile(fmt.Sprintf("%x", checksum), content)
	fmt.Printf("%x\n", checksum)
	return nil
}
