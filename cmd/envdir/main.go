package main

import (
	"fmt"
	"os"

	"github.com/mek/go-envdir"
)

func main() {
	// Check arguments
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: envdir dir child [args...]")
		os.Exit(111)
	}

	// Run child
	if err := envdir.Exec(os.Args[1], os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(111)
	}

}
