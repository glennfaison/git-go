package main

import (
	"flag"
	"fmt"
	"os"

	cat_file "github.com/codecrafters-io/git-starter-go/cmd/cat-file"
	hash_object "github.com/codecrafters-io/git-starter-go/cmd/hash-object"
	init_ "github.com/codecrafters-io/git-starter-go/cmd/init"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	flag.Parse()
	args := flag.Args()

	// fmt.Printf("%v\n", os.Args)
	// for i, arg := range args {
	// 	println(arg, i)
	// }

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	command, newArgs := args[0], args[1:]

	switch command {
	case "init":
		init_.CommandHandler_Init(newArgs)
	case "cat-file":
		cat_file.CommandHandler_CatFile(newArgs)
	case "hash-object":
		hash_object.CommandHandler_HashObject(newArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
