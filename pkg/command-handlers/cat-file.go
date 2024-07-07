package command_handlers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/git-starter-go/pkg"
)

func PrintCmd(blob_sha string) {
	blob_string, err := pkg.ReadObjectFile(blob_sha)
	pkg.CheckError(err)
	nullByteIndex := strings.Index(blob_string, "\x00")
	if nullByteIndex != -1 {
		blob_string = blob_string[nullByteIndex+1:]
	}
	fmt.Println(blob_string)
}

func CatFile(args []string) {
	flag := flag.NewFlagSet("mygit cat-file", flag.ExitOnError)
	blob_path := flag.String("p", "", "name of the blob file")
	flag.Parse(args)
	args = flag.Args()

	if blob_path != nil {
		PrintCmd(*blob_path)
	}
	if blob_path == nil {
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", args[0])
	}
}
