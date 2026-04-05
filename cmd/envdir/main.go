package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mek/go-envdir"
)

// main delegates to run and exits with the returned status code.
func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

// run implements the CLI behavior for envdir.
// It checks the argument shape, reports setup errors to stderr, and returns the
// exit code that main should hand back to the shell.
func run(args []string, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, "usage: envdir dir child [args...]")
		return 111
	}

	if err := envdir.Exec(args[0], args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 111
	}

	return 0
}
