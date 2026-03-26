# TODO

## Bugs

- Handle only regular files in `Read`.
  `readNames` currently accepts anything that is not a directory, which can include symlinks, FIFOs, sockets, and device nodes.
  This can lead to surprising reads or blocking on special files.
  Update the check to skip anything that is not a regular file.

## Improvements

- Clean up exported documentation in `envdir.go`.
  There are a few typos and awkward phrases in the package comments.

- Add a test for non-regular directory entries.
  This should cover the regular-file-only behavior once `Read` is tightened.

- Tighten `cmd/envdir` error output assertions.
  `TestRunReadError` currently checks only that stderr is non-empty.
  It should verify a stable, useful substring instead.
