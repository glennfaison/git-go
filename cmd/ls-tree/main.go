package ls_tree

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/pkg"
)

func CommandHandler_LsTree(args []string) {
	flag := flag.NewFlagSet("git ls-tree", flag.ExitOnError)
	nameOnly := flag.Bool("name-only", false, "prints only the file and folder names")
	flag.Parse(args)
	args = flag.Args()
	tree_sha := args[0]

	file_content, err := pkg.GetContentFromObject(tree_sha)
	pkg.CheckError(err)

	entries := ParseTreeObjectFromString(file_content)

	if nameOnly != nil {
		outputString := ""
		for _, entry := range entries {
			outputString += fmt.Sprintf("%s %s\t%s\n", entry.mode, entry.name, entry._20byteSha)
		}
		fmt.Print(outputString)
		return
	}

	outputString := ""
	for _, entry := range entries {
		outputString += fmt.Sprintf("%s\n", entry.name)
	}
	fmt.Print(outputString)
}

func ParseTreeObjectFromString(file_content string) []TreeObjectEntry {
	firstNullByteIndex := strings.Index(file_content, "\x00")
	body := ""
	if firstNullByteIndex > 0 {
		body = file_content[firstNullByteIndex+1:]
	}

	// RegExp for repeated sequences of `(040000 folder1\x007f21f4d392c2d79987c1)(100644 file1\x00d51d366274410103d3ec)...`
	bodyRegExp := regexp.MustCompile("(\\d{6}) (\\w+)\x00(\\w{20})")
	matches := bodyRegExp.FindAllStringSubmatch(body, -1)
	entries := []TreeObjectEntry{}
	for _, value := range matches {
		mode, name, _20byteSha := value[1], value[2], value[3]
		entries = append(entries, TreeObjectEntry{mode, name, _20byteSha})
	}
	return entries
}

type TreeObjectEntry struct {
	mode       string
	name       string
	_20byteSha string
}
