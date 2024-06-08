package cat_file

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/codecrafters-io/git-starter-go/cmd/pkg"
)

func P(blob_sha string) {
	blob_filename := path.Join(pkg.DOT_GIT_OBJECTS, blob_sha)
	blob_file, err := os.OpenFile(blob_filename, os.O_RDONLY, 0644)
	CheckError(err)
	zlibReader, err := zlib.NewReader(io.Reader(blob_file))
	CheckError(err)
	blob_bytes, err := io.ReadAll(zlibReader)
	CheckError(err)

	nullByteIndex := bytes.Index(blob_bytes, []byte("\x00"))
	blob_string := string(blob_bytes)
	if nullByteIndex != -1 {
		blob_string = string(blob_bytes[nullByteIndex:])
	}
	println(blob_string)
}

func CommandHandler_CatFile() {
	switch command := os.Args[2]; command {
	case "-p":
		blob_sha := os.Args[3]
		P(blob_sha)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
	}

}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
