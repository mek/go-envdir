// Package envdir reads environment variable settings from a directory.
//
// Each regular file in the directory becomes one environment entry. Empty
// files unset variables. Non-empty files contribute their first line after NUL
// bytes are converted to newlines and trailing spaces and tabs are removed.
//
// The package can read entries, apply them to an existing environment, or run
// a command with the resulting environment.
package envdir
