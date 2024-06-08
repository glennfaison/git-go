package main

import (
	"fmt"
	"os"

	cat_file "github.com/codecrafters-io/git-starter-go/cmd/cat-file"
	init_ "github.com/codecrafters-io/git-starter-go/cmd/init_"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		init_.CommandHandler_Init()
	case "cat-file":
		cat_file.CommandHandler_CatFile()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
