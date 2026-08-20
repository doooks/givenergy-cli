# Contributing

Contributions are welcome — bug reports, fixes, register-map corrections,
new metrics, whatever's useful.

A few small things before opening a PR:

- Run `go vet ./...` and `go test ./...`, and `gofmt -l .` should print
  nothing.
- If you're touching the protocol layer (`internal/givenergy/`), please
  test against a real dongle if you can, and say so in the PR.

For anything more than a small fix, opening an issue first to discuss the
approach is probably a good idea.
