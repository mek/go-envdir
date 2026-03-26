package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestRunUsage checks the CLI usage error path.
func TestRunUsage(t *testing.T) {
	var stderr bytes.Buffer

	code := run(nil, &stderr)
	if code != 111 {
		t.Fatalf("run() code = %d, want 111", code)
	}

	if stderr.String() != "usage: envdir dir child [args...]\n" {
		t.Fatalf("run() stderr = %q", stderr.String())
	}
}

// TestRunReadError checks that setup failures return the daemontools-style code.
func TestRunReadError(t *testing.T) {
	var stderr bytes.Buffer

	code := run([]string{filepath.Join(t.TempDir(), "missing"), "true"}, &stderr)
	if code != 111 {
		t.Fatalf("run() code = %d, want 111", code)
	}

	if stderr.Len() == 0 {
		t.Fatal("run() stderr is empty, want error output")
	}
}

// TestRunSuccess checks that the CLI hands the environment to a child process.
func TestRunSuccess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		if got := os.Getenv("TEST_VALUE"); got != "from-envdir" {
			fmt.Fprintf(os.Stderr, "TEST_VALUE = %q, want %q\n", got, "from-envdir")
			os.Exit(1)
		}
		os.Exit(0)
	}

	dir := t.TempDir()
	writeFile(t, dir, "TEST_VALUE", "from-envdir\n")

	var stderr bytes.Buffer
	args := []string{
		dir,
		os.Args[0],
		"-test.run=TestRunSuccess",
	}

	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	code := run(args, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, want 0, stderr = %q", code, stderr.String())
	}
}

// writeFile writes one test file for the CLI tests.
func writeFile(t *testing.T, dir, name, contents string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}
