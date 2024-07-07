package command_handlers

import (
	"flag"
	"fmt"

	"github.com/codecrafters-io/git-starter-go/pkg"
)

func CommandHandler_CommitTree(args []string) error {
	flag := flag.NewFlagSet("mygit commit-tree", flag.ExitOnError)
	parents := flag.String("p", "", "id of a parent commit object")
	message := flag.String("m", "", "commit message")
	flag.Parse(args)
	treeObject := args[0]

	checksum, content, err := pkg.ComputeCommitObject(treeObject, []string{*parents}, *message)
	if err != nil {
		return err
	}

	pkg.WriteObjectFile(fmt.Sprintf("%x", checksum), content)
	fmt.Printf("%x\n", checksum)
	return nil
}
