package hash_object

import (
	"flag"
	"fmt"
)

func CommandHandler_HashObject(args []string) {
	flag := flag.NewFlagSet("mygit hash-object", flag.ExitOnError)
	write := flag.Bool("w", false, "creates a git object for the file in .git/objects")
	flag.Parse(args)
	args = flag.Args()

	fmt.Printf("%v\n", args)
	fmt.Println(write)

	// objectText := GetObjectText()
	// if write != nil {
	// }
}
