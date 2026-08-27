# Consuming go-srd5e

`go-srd5e` is a public Go module. Add it like any other dependency:

```bash
go get github.com/trancecode/go-srd5e@latest
```

Browse the API docs at <https://pkg.go.dev/github.com/trancecode/go-srd5e>.

For day-to-day development across several repos at once, point a consumer's
`go.mod` at a local checkout so changes don't require tagging on every edit:

```
replace github.com/trancecode/go-srd5e => ../go-srd5e
```

Consumers pin a tagged version (for example `v1.2.0`); the `replace` directive
above is for local, cross-repo development only and should not land in a
committed `go.mod`.
