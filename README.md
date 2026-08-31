# ro.am/roamhq

Official Go SDK for the [Roam API](https://developer.ro.am).

```bash
go get ro.am/roamhq
```

Requires **Go 1.22+**. The module has no dependencies beyond the Go standard
library (the generated client may add small ones Fern needs).

```go
package main

import (
	"context"
	"os"

	roamhq "ro.am/roamhq"
	"ro.am/roamhq/client"
	"ro.am/roamhq/option"
)

func main() {
	c := client.NewClient(option.WithToken(os.Getenv("ROAM_TOKEN")))

	_, err := c.Chat.Post(context.Background(), &roamhq.PostChatRequest{
		GroupId: roamhq.String("88bebce7-6cbb-4666-96f9-5c02d73e6661"),
		Text:    roamhq.String("Build completed successfully!"),
	})
	if err != nil {
		panic(err)
	}
}
```

Full usage — pagination, retries, error handling, version pinning — is in the
[SDK guide](https://developer.ro.am/docs/guides/sdks). The
[API reference](https://developer.ro.am/docs/api/api) documents every endpoint.

## Verifying webhooks

```go
import "ro.am/roamhq/webhooks"

event, err := webhooks.Verify(body, r.Header, os.Getenv("ROAM_WEBHOOK_SECRET"), nil)
if err != nil {
    http.Error(w, "unauthorized", http.StatusUnauthorized)
    return
}
```

`body` must be the **raw request bytes**. Treat any error as a `401`. Pass the
signing secret exactly as Roam issued it, `whsec_` prefix included.

This is the one part of the SDK that is hand-written rather than generated. See
[`webhooks/verify.go`](webhooks/verify.go) for why.

## This repository is mostly generated

The client, types, and option packages are generated from the Roam OpenAPI
specification by [Fern](https://buildwithfern.com). **Do not edit them by
hand** — the next regeneration will overwrite them.

A change to the spec opens or updates a single long-lived pull request here
titled "Regenerate SDK from OpenAPI spec". A human reviews the diff, picks the
semantic-version bump (the git tag), and merges.

Everything listed in [`.fernignore`](.fernignore) is hand-maintained:

| Path | What it is |
| --- | --- |
| `webhooks/` | Hand-written signature verification, imported as `ro.am/roamhq/webhooks`. |
| `tests/` | Contract tests. See [`tests/README.md`](tests/README.md). |
| `.github/` | CI. |
| `README.md`, `LICENSE` | This file, and the license. |

If you find a wrong type or a missing field, it is almost always a bug in the
[OpenAPI spec](https://developer.ro.am/chat-v1.json) rather than in hand-written
code — fixing it there fixes the docs site and every SDK at once. Report it via
[Roam Support](https://ro.am/support/contact-us) or
[developer@ro.am](mailto:developer@ro.am).

## Building

```bash
go build ./...
go vet ./...
go test ./webhooks/... ./tests/...
```

Do not `go test ./...`. Fern emits WireMock-backed tests that expect a mock
server on localhost:8080; they are scaffolding, not assertions about this SDK.
