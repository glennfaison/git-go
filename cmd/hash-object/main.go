package hash_object

import (
	"compress/zlib"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/pkg"
)

func CommandHandler_HashObject(args []string) {
	flag := flag.NewFlagSet("mygit hash-object", flag.ExitOnError)
	write := flag.Bool("w", false, "creates a git object for the file in .git/objects")
	flag.Parse(args)
	args = flag.Args()

	filePath := args[0]
	fileContents, err := os.ReadFile(filePath)
	pkg.CheckError(err)

	fileContentString := string(fileContents)
	shaInput := "blob " + strconv.Itoa(len(fileContents)) + "\x00" + fileContentString

	shaOutput := sha1.Sum([]byte(shaInput))
	hash := fmt.Sprintf("%x", shaOutput)

	// Perform the key requirement of the `hash-object` command: print the SHA.
	fmt.Print(hash)

	if write != nil {
		WriteToObjects(hash, shaInput)
	}
}

func WriteToObjects(filename string, contentString string) {
	// Calculate the paths to the file and its parent directory.
	filePath := path.Join(pkg.DOT_GIT_OBJECTS, filename[:2], filename[2:])
	dirPath := strings.TrimSuffix(filePath, filename[2:])

	// Make the parent directory structure if necessary.
	err := os.MkdirAll(dirPath, os.ModePerm)
	pkg.CheckError(err)

	// Write to a temporary file, and only replace the original after a success.
	f, err := os.CreateTemp(dirPath, "*")
	pkg.CheckError(err)
	defer f.Close()
	defer func() {
		err := os.Remove(path.Join(dirPath, f.Name()))
		pkg.CheckNonIsNotExistError(err)
	}()

	zlibWriter := zlib.NewWriter(io.Writer(f))
	_, err = zlibWriter.Write([]byte(contentString))
	pkg.CheckError(err)
	defer zlibWriter.Close()

	err = os.Remove(filePath)
	pkg.CheckNonIsNotExistError(err)
	os.Rename(f.Name(), filePath)
}
