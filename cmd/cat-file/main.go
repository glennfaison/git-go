package cat_file

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/codecrafters-io/git-starter-go/cmd/pkg"
)

func PrintCmd(blob_sha string) {
	// cmd := exec.Command("ls", "-l", ".")
	// stdout, err := cmd.Output()

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// // Print the output
	// fmt.Println("pwd ->", string(stdout))

	blob_filename := path.Join(pkg.DOT_GIT_OBJECTS, blob_sha[:2], blob_sha[2:])
	println("before OpenFile")
	blob_file, err := os.OpenFile(blob_filename, os.O_RDONLY, 0644)
	CheckError(err)
	defer blob_file.Close()
	println("after OpenFile")
	zlibReader, err := zlib.NewReader(io.Reader(blob_file))
	CheckError(err)
	defer zlibReader.Close()
	blob_bytes, err := io.ReadAll(zlibReader)
	CheckError(err)

	nullByteIndex := bytes.Index(blob_bytes, []byte("0"))
	blob_string := string(blob_bytes)
	fmt.Println("blob_string:", blob_string)
	if nullByteIndex != -1 {
		blob_string = string(blob_bytes[nullByteIndex:])
	}
	println(blob_string)
}

func CommandHandler_CatFile(args []string) {
	// fmt.Printf("%v\n", args)
	flag := flag.NewFlagSet("mygit cat-file", flag.ExitOnError)
	blob_path := flag.String("p", "", "name of the blob file")
	flag.Parse(args)
	args = flag.Args()

	switch {
	case blob_path != nil:
		PrintCmd(*blob_path)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", args[0])
	}

}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", err.Error())
		os.Exit(0)
	}
}
