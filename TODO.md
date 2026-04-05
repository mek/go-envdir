# TODO

## Improvements

- Clean up exported documentation in `envdir.go`.
  There are a few typos and awkward phrases in the package comments.

- Tighten `cmd/envdir` error output assertions.
  `TestRunReadError` currently checks only that stderr is non-empty.
  It should verify a stable, useful substring instead.
