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

Drop the `replace` and pin a tagged version (for example `v0.1.0`) once the
module is stable.
