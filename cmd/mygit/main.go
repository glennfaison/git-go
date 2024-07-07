package main

import (
	"flag"
	"fmt"
	"os"

	command_handlers "github.com/codecrafters-io/git-starter-go/pkg/command-handlers"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	command, newArgs := args[0], args[1:]

	switch command {
	case "init":
		command_handlers.Init(newArgs)
	case "cat-file":
		command_handlers.CatFile(newArgs)
	case "hash-object":
		command_handlers.HashObject(newArgs)
	case "ls-tree":
		command_handlers.LsTree(newArgs)
	case "write-tree":
		command_handlers.WriteTree(newArgs)
	case "commit-tree":
		command_handlers.CommitTree(newArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
