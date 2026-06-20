# Consuming go-srd5e

This module is private under the trancecode organization.

```bash
go env -w GOPRIVATE=github.com/trancecode/*
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

For local development across games, add to the consuming game's `go.mod`:

```
replace github.com/trancecode/go-srd5e => ../go-srd5e
```

Drop the replace and pin a tagged version once the module is stable.
