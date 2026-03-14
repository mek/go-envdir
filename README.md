# envdir (Go)

A small Go implementation of **envdir**, originally from D. J. Bernstein's `daemontools`.

`envdir` reads environment variables from files in a directory and applies them when running a command.

It is useful for container environments, secrets stored as files, and simple configuration setups.

This project provides:

* a small reusable **Go library**
* a minimal **CLI tool**

The implementation favors **simplicity, clarity, and Unix-style design**.

---

# Installation

Install the CLI using Go:

```
go install github.com/mek/go-envdir/cmd/envdir@latest
```

The binary will be placed in your `$GOBIN` or `$GOPATH/bin`.

---

# CLI Usage

```
envdir DIR COMMAND [ARGS...]
```

Example:

```
envdir ./env ./myapp
```

This will:

1. Read environment variables from `./env`
2. Apply them to the current environment
3. Run `./myapp`

---

# Directory Format

Each file in the directory represents one environment variable.

Example:

```
env/
  DB_HOST
  DB_PORT
  DEBUG
```

### File name

The filename becomes the variable name.

```
env/DB_HOST
```

### File contents

The **first line** becomes the value.

```
localhost
```

Result:

```
DB_HOST=localhost
```

---

# Rules

Behavior follows the traditional `envdir` design.

| Condition            | Result                   |
| -------------------- | ------------------------ |
| File empty           | variable is **unset**    |
| File contains text   | first line becomes value |
| Trailing spaces/tabs | removed                  |
| NUL bytes            | converted to newline     |

---

# Example

Directory:

```
env/
  DB_HOST
  DB_PORT
  DEBUG
```

Contents:

```
env/DB_HOST → localhost
env/DB_PORT → 5432
env/DEBUG   → true
```

Command:

```
envdir ./env ./server
```

Environment seen by `server`:

```
DB_HOST=localhost
DB_PORT=5432
DEBUG=true
```

---

# Library Usage

The Go package can also be used directly.

```go
package main

import (
	"fmt"
	"os"

	"github.com/mek/go-envdir"
)

func main() {
	entries, err := envdir.Read("./env")
	if err != nil {
		panic(err)
	}

	env := envdir.Apply(os.Environ(), entries)

	for _, e := range env {
		fmt.Println(e)
	}
}
```

---

# Design Philosophy

This project intentionally stays small.

Goals:

* simple code
* minimal API
* standard library only
* predictable behavior

Inspired by the design approach of **Rob Pike**, **Brian Kernighan**, and classic Unix tools.

---

# Status

Early development.

The core functionality is implemented and the API may still evolve.

---

# License

MIT

