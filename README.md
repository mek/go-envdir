# envdir (Go)

A small Go implementation of **envdir**, originally from D. J. Bernstein's `daemontools`.

`envdir` reads environment variables from files in a directory and applies them when running a command.

It is useful for container environments, secrets stored as files, and simple configuration setups.

This project provides:

* a small reusable **Go library**
* a minimal **CLI tool**
* a small **sample file server**

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

# Sample File Server

There is a small sample program in `examples/fileserver`.

It reads an envdir directory, applies the entries to the current environment,
and then serves files with `http.FileServer`.

Build it with:

```sh
make build-examples
```

That produces:

```text
./bin/fileserver
```

Environment variables used by the sample:

* `SERVER_HOST` default `127.0.0.1`
* `SERVER_PORT` default `8080`
* `SERVER_DIR` default `.`

Example envdir layout:

```
examples/fileserver/env/
  SERVER_HOST
  SERVER_PORT
  SERVER_DIR
```

Run it directly:

```sh
go run ./examples/fileserver ./examples/fileserver/env
```

Or run the built binary:

```sh
./bin/fileserver ./examples/fileserver/env
```

Then visit:

```text
http://127.0.0.1:8090/
```

You can also run the same program through the CLI:

```sh
go run ./cmd/envdir ./examples/fileserver/env go run ./examples/fileserver
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
