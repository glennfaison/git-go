package pkg

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
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

func WriteObjectFile(filename string, content []byte) {
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
	_, err = zlibWriter.Write(content)
	CheckError(err)
	defer zlibWriter.Close()

	err = os.Remove(filePath)
	CheckNonIsNotExistError(err)
	os.Rename(f.Name(), filePath)
}

func ReadObjectFile(filename string) (string, error) {
	blob_filename := path.Join(DOT_GIT_OBJECTS, filename[:2], filename[2:])
	blob_file, err := os.OpenFile(blob_filename, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return "", err
	}
	defer blob_file.Close()

	zlibReader, err := zlib.NewReader(io.Reader(blob_file))
	if err != nil {
		return "", err
	}
	defer zlibReader.Close()
	blob_bytes, err := io.ReadAll(zlibReader)
	if err != nil {
		err := fmt.Errorf("error while reading file: %s\n%w", blob_file.Name(), err)
		return "", err
	}

	blob_string := string(blob_bytes)
	return blob_string, nil
}

func ComputeBlobObjectForFile(filePath string) ([20]byte, []byte, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		err = fmt.Errorf("failed to open file: %s\n%w", filePath, err)
		return [20]byte{}, []byte{}, err
	}
	fileContents, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("failed to read file: %s\n%w", filePath, err)
		return [20]byte{}, []byte{}, err
	}
	header := fmt.Sprintf("blob %d\x00", len(fileContents))
	sha1Input := append([]byte(header), fileContents...)
	return sha1.Sum(sha1Input), sha1Input, nil
}

func ParseTreeObjectFromString(file_content string) []TreeObjectEntry {
	firstNullByteIndex := strings.Index(file_content, "\x00")
	body := ""
	if firstNullByteIndex > 0 {
		body = file_content[firstNullByteIndex+1:]
	}

	// RegExp for repeated sequences of `(040000 folder1\x007f21f4d392c2d79987c1)(100644 file1\x00d51d366274410103d3ec)...`
	bodyRegExp := regexp.MustCompile("(\\d{6}) ([^\x00]+)\x00(\\w{20})")
	matches := bodyRegExp.FindAllStringSubmatch(body, -1)
	entries := []TreeObjectEntry{}
	for _, value := range matches {
		mode, name, _20byteSha := value[1], value[2], value[3]
		entries = append(entries, TreeObjectEntry{mode, name, _20byteSha})
	}
	return entries
}

type TreeObjectEntry struct {
	Mode         string
	Name         string
	ShaAs20Bytes string
}

func ComputeTreeObjectForDirectory(dir string, writeToFile bool) ([20]byte, []byte, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return [20]byte{}, nil, err
	}

	treeObjEntries := []TreeObjectEntry{}
	for _, dirEntry := range dirEntries {
		filename := path.Join(dir, dirEntry.Name())
		if filename == ".git" {
			continue
		}
		if dirEntry.IsDir() {
			checksum, _, err := ComputeTreeObjectForDirectory(filename, writeToFile)
			if err != nil {
				return [20]byte{}, nil, err
			}

			item := TreeObjectEntry{"040000", dirEntry.Name(), string(checksum[:])}
			treeObjEntries = append(treeObjEntries, item)
			continue
		}
		checksum, _, err := ComputeBlobObjectForFile(filename)
		if err != nil {
			return [20]byte{}, nil, err
		}

		item := TreeObjectEntry{"100644", dirEntry.Name(), string(checksum[:])}
		treeObjEntries = append(treeObjEntries, item)
	}

	sort.Slice(treeObjEntries, func(i, j int) bool {
		return treeObjEntries[i].Name < treeObjEntries[j].Name
	})
	body := []byte{}
	for _, treeObjEntry := range treeObjEntries {
		body = append(body, treeObjEntry.ShaAs20Bytes...)
	}

	header := fmt.Sprintf("tree %d\x00", len(body))
	body = append([]byte(header), body...)

	checksum := sha1.Sum(body)

	return checksum, body, nil
}
