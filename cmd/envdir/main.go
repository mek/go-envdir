package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mek/go-envdir"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

func run(args []string, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, "usage: envdir dir child [args...]")
		return 111
	}

	if err := envdir.Exec(args[0], args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		return 111
	}

	return 0
}
