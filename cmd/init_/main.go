package init_

import (
	"fmt"
	"os"

	pkg "github.com/codecrafters-io/git-starter-go/cmd/pkg"
)

func CommandHandler_Init() {
	for _, dir := range []string{pkg.DOT_GIT, pkg.DOT_GIT_OBJECTS, pkg.DOT_GIT_REFS} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(pkg.DOT_GIT_HEAD, headFileContents, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
	}

	fmt.Println("Initialized git directory")
}
