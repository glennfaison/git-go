package pkg

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

const DOT_GIT = ".git"
const DOT_GIT_OBJECTS = ".git/objects"
const DOT_GIT_REFS = ".git/refs"
const DOT_GIT_HEAD = ".git/HEAD"

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(0)
	}
}

func CheckNonIsNotExistError(err error) {
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(0)
	}
}

func WriteToObjects(filename string, contentString string) {
	// Calculate the paths to the file and its parent directory.
	filePath := path.Join(DOT_GIT_OBJECTS, filename[:2], filename[2:])
	dirPath := strings.TrimSuffix(filePath, filename[2:])

	// Make the parent directory structure if necessary.
	err := os.MkdirAll(dirPath, os.ModePerm)
	CheckError(err)

	// Write to a temporary file, and only replace the original after a success.
	f, err := os.CreateTemp(dirPath, "*")
	CheckError(err)
	defer f.Close()
	defer func() {
		err := os.Remove(path.Join(dirPath, f.Name()))
		CheckNonIsNotExistError(err)
	}()

	zlibWriter := zlib.NewWriter(io.Writer(f))
	_, err = zlibWriter.Write([]byte(contentString))
	CheckError(err)
	defer zlibWriter.Close()

	err = os.Remove(filePath)
	CheckNonIsNotExistError(err)
	os.Rename(f.Name(), filePath)
}

func GetContentFromObject(filename string) (string, error) {
	blob_filename := path.Join(DOT_GIT_OBJECTS, filename[:2], filename[2:])
	blob_file, err := os.OpenFile(blob_filename, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return "", err
	}
	defer blob_file.Close()

	zlibReader, err := zlib.NewReader(io.Reader(blob_file))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return "", err
	}
	defer zlibReader.Close()
	blob_bytes, err := io.ReadAll(zlibReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return "", err
	}

	blob_string := string(blob_bytes)
	return blob_string, nil
}
