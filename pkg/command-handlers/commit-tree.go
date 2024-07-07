package command_handlers

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/codecrafters-io/git-starter-go/pkg"
)

type StringSlice []string

func (s *StringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var _ flag.Value = (*StringSlice)(nil)

func CommitTree(args []string) error {
	treeObjectSha := args[0]
	if matched, _ := regexp.MatchString(`^[0-9a-fA-F]+$`, treeObjectSha); !matched {
		return fmt.Errorf("expected to receive a valid sha")
	}
	args = args[1:]

	flag := flag.NewFlagSet("mygit commit-tree", flag.ExitOnError)
	parent := flag.String("p", "", "id of parent commit object")
	message := flag.String("m", "", "commit message")
	flag.Parse(args)

	checksum, content, err := pkg.ComputeCommitObject(treeObjectSha, *parent, *message)
	if err != nil {
		return err
	}

	pkg.WriteObjectFile(fmt.Sprintf("%x", checksum), content)
	fmt.Printf("%x\n", checksum)
	return nil
}
