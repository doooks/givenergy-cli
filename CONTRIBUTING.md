# Contributing

Contributions are welcome — bug reports, fixes, register-map corrections,
new metrics, whatever's useful.

A few small things before opening a PR:

- Run `go vet ./...` and `go test ./...`, and `gofmt -l .` should print
  nothing.
- If you're touching the protocol layer (`internal/givenergy/`), please
  test against a real dongle if you can, and say so in the PR — sign
  conventions and register behavior have bitten this project before (see
  the commit history) and are hard to catch from reading the code alone.
- Keep changes focused; this is meant to stay a small, lightweight tool.

For anything more than a small fix, opening an issue first to discuss the
approach is appreciated but not required.
