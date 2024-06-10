package hash_object

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"os"
	"strconv"

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
	fmt.Printf("%s", hash)

	if write != nil {
		pkg.WriteToObjects(hash, shaInput)
	}
}
