# AGENTS.md

This is the generated Go SDK for the Roam API, published as `ro.am/roamhq`.

## What you may edit

Hand-written files are listed in `.fernignore`. Those are the only files a
regeneration will not overwrite:

- `webhooks/` — signature verification. Security-relevant; tests live next to
  the source. Do not replace this with the generated helper.
- `README.md`, `LICENSE`, `RELEASING.md`, `AGENTS.md`
- `.github/`, `.gitignore`, `.fernignore`

## What you must not edit

Everything else — `client/`, `option/`, `core/`, `internal/`, root `*.go`,
`go.mod` — is generated from
[WonderInventions/developer-ro-am](https://github.com/WonderInventions/developer-ro-am)
by Fern. A spec change opens a PR titled "Regenerate SDK from OpenAPI spec"
that replaces this tree with `rsync --delete`.

A wrong type or missing field is almost always a bug in the OpenAPI spec, not
in this repo.

## Tests

```bash
go build ./...
go vet ./...
go test ./webhooks/...
```

Never `go test ./...`. The generator emits WireMock tests that dial
localhost:8080 and fail without a harness we do not run.
