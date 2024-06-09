package cat_file

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/pkg"
)

func PrintCmd(blob_sha string) {
	blob_filename := path.Join(pkg.DOT_GIT_OBJECTS, blob_sha[:2], blob_sha[2:])
	blob_file, err := os.OpenFile(blob_filename, os.O_RDONLY, 0644)
	CheckError(err)
	defer blob_file.Close()

	zlibReader, err := zlib.NewReader(io.Reader(blob_file))
	CheckError(err)
	defer zlibReader.Close()
	blob_bytes, err := io.ReadAll(zlibReader)
	CheckError(err)

	blob_string := string(blob_bytes)
	nullByteIndex := strings.Index(blob_string, "\x00")
	if nullByteIndex != -1 {
		blob_string = blob_string[nullByteIndex+1:]
	}
	fmt.Print(blob_string)
}

func CommandHandler_CatFile(args []string) {
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

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", err.Error())
		os.Exit(0)
	}
}
