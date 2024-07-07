package command_handlers

import (
	"flag"
	"fmt"

	"github.com/codecrafters-io/git-starter-go/pkg"
)

func HashObject(args []string) error {
	flag := flag.NewFlagSet("mygit hash-object", flag.ExitOnError)
	write := flag.Bool("w", false, "creates a git object for the file in .git/objects")
	flag.Parse(args)
	args = flag.Args()

	filePath := args[0]
	checksum, content, err := pkg.ComputeBlobObjectForFile(filePath)
	if err != nil {
		return err
	}
	hash := fmt.Sprintf("%x", checksum)

	// Perform the key requirement of the `hash-object` command: print the SHA.
	fmt.Printf("%s", hash)

	if write != nil && *write {
		pkg.WriteObjectFile(hash, content)
	}
	return nil
}
