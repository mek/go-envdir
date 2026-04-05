package envdir

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Entry describes one environment change read from a file.
// A file may contribute a value, or it may quietly insist that the variable
// should not exist at all.
type Entry struct {
	Name  string
	Value string
	Unset bool
}

// Read reads dir and returns the environment changes described by its files.
//
// Each file in dir becomes one entry. The file name is the environment
// variable name, which is an arrangement simple enough to be trusted.
//
// If the file is empty, the variable is unset.
//
// If the file is not empty, the value comes from the first line.
// Trailing spaces and tabs are removed, because they rarely improve matters.
//
// NUL bytes are turned into newlines before the first line is chosen.
func Read(dir string) ([]Entry, error) {
	names, err := readNames(dir)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(names))
	for _, name := range names {
		entry, err := readEntry(dir, name)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Apply applies entries to base and returns the resulting environment.
//
// Entries add, replace, or remove variables. The final result is sorted by
// name so that the outcome is stable even when the wider world is not.
func Apply(base []string, entries []Entry) []string {

	// Let's create an array for the current base
	env := make(map[string]string, len(base))
	for _, s := range base {
		name, value, ok := strings.Cut(s, "=")
		if !ok {
			continue
		}
		env[name] = value
	}

	// Unset needed entries from current env
	// Set new values as needed.
	for _, e := range entries {
		if e.Unset {
			delete(env, e.Name)
			continue
		}
		env[e.Name] = e.Value
	}

	// Get a list of the current variable names
	// and sort them
	names := make([]string, 0, len(env))
	for name := range env {
		names = append(names, name)
	}
	sort.Strings(names)

	// Create a new env array for output
	// using the requested entries.
	out := make([]string, 0, len(names))
	for _, name := range names {
		out = append(out, name+"="+env[name])
	}
	return out
}

// Exec runs argv[0] with argv[1:] using environment settings read from dir.
// The current process environment is used as the base before those changes are
// applied.
//
// If argv is empty, Exec returns an error rather than pretending a command
// might appear if everyone waits long enough.
func Exec(dir string, argv []string) error {

	// bravely refuse to do nothing
	if len(argv) == 0 {
		return fmt.Errorf("exec: no command")
	}

	// read the entires
	entries, err := Read(dir)
	if err != nil {
		return err
	}

	// Let's setup the child to run
	cmd := exec.Command(argv[0], argv[1:]...)
	// Let's get the new environmanet applied to the command
	cmd.Env = Apply(os.Environ(), entries)
	// Output the the current process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Let's do this, giving this process the return code of the
	// command we run.
	return cmd.Run()
}

// readNames returns the sorted variable names found in dir.
// Only regular files are considered. Directories, symlinks, and other special
// entries are ignored before variable-name validation.
func readNames(dir string) ([]string, error) {
	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}

	names := make([]string, 0, len(items))
	for _, item := range items {
		info, err := item.Info()
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", filepath.Join(dir, item.Name()), err)
		}
		if !info.Mode().IsRegular() {
			continue
		}

		name := item.Name()
		if !validName(name) {
			return nil, fmt.Errorf("invalid variable name %q", name)
		}

		names = append(names, name)
	}

	sort.Strings(names)
	return names, nil
}

// readEntry reads one file and turns it into a single environment change.
// Empty files unset variables; non-empty files contribute their first line
// after NUL conversion and trailing space trimming.
func readEntry(dir, name string) (Entry, error) {
	path := filepath.Join(dir, name)

	data, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, fmt.Errorf("read %s: %w", path, err)
	}

	if len(data) == 0 {
		return Entry{
			Name:  name,
			Unset: true,
		}, nil
	}

	data = bytes.ReplaceAll(data, []byte{0}, []byte{'\n'})

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		data = data[:i]
	}

	value := strings.TrimRight(string(data), " \t")

	return Entry{
		Name:  name,
		Value: value,
	}, nil
}

// validName reports whether name can stand as an environment variable name.
// The rules are small: no empty names, no '=', no '/', and no NUL bytes.
func validName(name string) bool {
	if name == "" {
		return false
	}
	if strings.ContainsRune(name, '=') {
		return false
	}
	if strings.ContainsRune(name, '/') {
		return false
	}
	if strings.ContainsRune(name, 0) {
		return false
	}
	return true
}
