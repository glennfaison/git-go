package pkg

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DOT_GIT = ".git"
const DOT_GIT_OBJECTS = ".git/objects"
const DOT_GIT_REFS = ".git/refs"
const DOT_GIT_HEAD = ".git/HEAD"

type TreeObjectEntry struct {
	Mode         string
	Name         string
	ShaAs20Bytes string
	Type         string
}

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
	byteSlice := []byte(file_content)
	startIndex, endIndex := 0, 0

	// Read the "tree" prefix.
	endIndex = bytes.IndexByte(byteSlice, byte(' '))
	if string(byteSlice[startIndex:endIndex]) != "tree" {
		return nil
	}
	startIndex = endIndex + 1

	// Read the file size.
	endIndex = startIndex + bytes.IndexByte(byteSlice[startIndex:], byte(0))
	_, err := strconv.Atoi(string(byteSlice[startIndex:endIndex]))
	if err != nil {
		return nil
	}
	startIndex = endIndex + 1

	entries := []TreeObjectEntry{}
	length := len(byteSlice[startIndex:])
	for endIndex <= length {
		// Read the Mode (a number).
		endIndex = startIndex + bytes.IndexByte(byteSlice[startIndex:], byte(' '))
		mode := string(byteSlice[startIndex:endIndex])
		mode = fmt.Sprintf("%06s", mode)
		_, err := strconv.Atoi(mode)
		if err != nil {
			return nil
		}
		startIndex = endIndex + 1

		// Read the Name (a string).
		endIndex = startIndex + bytes.IndexByte(byteSlice[startIndex:], byte(0))
		name := string(byteSlice[startIndex:endIndex])
		startIndex = endIndex + 1

		// Read the 20-byte-long hash (a string).
		endIndex = startIndex + 20
		hashString := string(byteSlice[startIndex:endIndex])
		startIndex = endIndex

		// Append the entry we just read.
		objectType := "blob"
		if mode == "040000" {
			objectType = "tree"
		}
		entries = append(entries, TreeObjectEntry{Mode: mode, Type: objectType, Name: name, ShaAs20Bytes: hashString})
	}
	return entries
}

func ComputeTreeObjectForDirectory(dir string, writeToFile bool) ([20]byte, []byte, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return [20]byte{}, nil, err
	}

	ignoredDirectories := []string{".git"}
	treeObjEntries := []TreeObjectEntry{}
	for _, dirEntry := range dirEntries {
		filename := path.Join(dir, dirEntry.Name())
		if slices.Contains(ignoredDirectories, filename) {
			continue
		}

		var checksum [20]byte
		var content []byte
		var err error
		var item TreeObjectEntry = TreeObjectEntry{}
		if dirEntry.IsDir() {
			item = TreeObjectEntry{Mode: "40000", Type: "tree"}
			checksum, content, err = ComputeTreeObjectForDirectory(filename, writeToFile)
		} else {
			item = TreeObjectEntry{Mode: "100644", Type: "blob"}
			checksum, content, err = ComputeBlobObjectForFile(filename)
		}
		if err != nil {
			return [20]byte{}, nil, err
		}
		item.Name = dirEntry.Name()
		item.ShaAs20Bytes = string(checksum[:])

		if writeToFile {
			WriteObjectFile(fmt.Sprintf("%x", checksum), content)
		}
		treeObjEntries = append(treeObjEntries, item)
	}

	sort.Slice(treeObjEntries, func(i, j int) bool {
		return treeObjEntries[i].Name < treeObjEntries[j].Name
	})

	body := []byte{}
	for _, treeObjEntry := range treeObjEntries {
		entryTemplate := fmt.Sprintf("%s %s\x00%s", treeObjEntry.Mode, treeObjEntry.Name, treeObjEntry.ShaAs20Bytes)
		body = append(body, []byte(entryTemplate)...)
	}

	header := fmt.Sprintf("tree %d\x00", len(body))
	body = append([]byte(header), body...)

	checksum := sha1.Sum(body)

	return checksum, body, nil
}

func ComputeCommitObject(treeObjectSha string, parentShas []string, message string) ([20]byte, []byte, error) {
	userName := "Anonymous Developer"
	userEmail := "anonymous@example.com"
	committerName := "Anonymous Developer"
	committerEmail := "anonymous@example.com"

	bodyString := fmt.Sprintf(
		"tree %s\nparent %s\nauthor %s <%s> %d %s\ncommitter %s <%s> %d %s\n\n%s",
		treeObjectSha, parentShas[0],
		userName, userEmail, time.Now().Unix(), time.Now().Format(time.RFC3339),
		committerName, committerEmail, time.Now().Unix(), time.Now().Format(time.RFC3339), message,
	)
	body := []byte(bodyString)
	header := fmt.Sprintf("commit %d\x00", len(body))

	content := []byte(header)
	content = append(content, []byte(body)...)
	checksum := sha1.Sum(content)

	return checksum, content, nil
}
