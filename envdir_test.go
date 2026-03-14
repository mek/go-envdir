package envdir

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "EMPTY", "")
	writeFile(t, dir, "FIRST", "value\nsecond\n")
	writeFile(t, dir, "TRIM", "value \t\t")
	writeFile(t, dir, "NUL", "left\x00right")
	writeFile(t, dir, "ALPHA", "a")

	entries, err := Read(dir)
	if err != nil {
		t.Fatalf("Read(%q) error = %v", dir, err)
	}

	want := []Entry{
		{Name: "ALPHA", Value: "a"},
		{Name: "EMPTY", Unset: true},
		{Name: "FIRST", Value: "value"},
		{Name: "NUL", Value: "left"},
		{Name: "TRIM", Value: "value"},
	}

	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("Read(%q) = %#v, want %#v", dir, entries, want)
	}
}

func TestReadInvalidName(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "GOOD", "ok")
	writeFile(t, dir, "BAD=NAME", "bad")

	_, err := Read(dir)
	if err == nil {
		t.Fatal("Read() error = nil, want invalid variable name error")
	}
	if !strings.Contains(err.Error(), `invalid variable name "BAD=NAME"`) {
		t.Fatalf("Read() error = %q, want invalid variable name", err)
	}
}

func TestApply(t *testing.T) {
	base := []string{
		"KEEP=1",
		"REPLACE=old",
		"REMOVE=gone",
		"INVALID",
	}

	entries := []Entry{
		{Name: "REPLACE", Value: "new"},
		{Name: "ADD", Value: "2"},
		{Name: "REMOVE", Unset: true},
	}

	got := Apply(base, entries)
	want := []string{
		"ADD=2",
		"KEEP=1",
		"REPLACE=new",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Apply() = %#v, want %#v", got, want)
	}
}

func TestExecNoCommand(t *testing.T) {
	err := Exec(t.TempDir(), nil)
	if err == nil {
		t.Fatal("Exec() error = nil, want error")
	}
	if err.Error() != "exec: no command" {
		t.Fatalf("Exec() error = %q, want %q", err, "exec: no command")
	}
}

func writeFile(t *testing.T, dir, name, contents string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}
