package ls_tree

import (
	"flag"
	"fmt"

	"github.com/codecrafters-io/git-starter-go/pkg"
)

func CommandHandler_LsTree(args []string) {
	flag := flag.NewFlagSet("git ls-tree", flag.ExitOnError)
	nameOnly := flag.Bool("name-only", false, "prints only the file and folder names")
	flag.Parse(args)
	args = flag.Args()
	tree_sha := args[0]

	file_content, err := pkg.ReadObjectFile(tree_sha)
	pkg.CheckError(err)

	entries := pkg.ParseTreeObjectFromString(file_content)
	fmt.Printf("\"GOT FILE CONTENT!\": %v\n", file_content)

	if nameOnly != nil {
		outputString := ""
		for _, entry := range entries {
			outputString += fmt.Sprintf("%s %s\t%s\n", entry.Mode, entry.Name, entry.ShaAs20Bytes)
		}
		fmt.Printf("outputString: %v\n", outputString)
		return
	}

	outputString := ""
	for _, entry := range entries {
		outputString += fmt.Sprintf("%s\n", entry.Name)
	}
	fmt.Print(outputString)
}
