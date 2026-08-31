# Contract tests

These tests exist because the generated SDK has no other behavioral safety net.

Fern can emit WireMock wire tests, but they expect a mock server on
localhost:8080 and are not run in CI. Everything upstream of this repo checks
*shapes*: `generate-sdks.yml` in developer-ro-am builds and vets the generated
sources. Nothing in that chain proves the client *behaves* the way the spec
says it should.

That gap matters most for the cross-cutting behavior modeled in
`fern/apis/roam/overrides.yml` and `overlays.yml`. A mistake there produces an
SDK that compiles perfectly and is quietly wrong:

- A pagination field renamed in the spec makes the client return page one and
  report no next page. That reads as "the account only has two groups".
- A `Roam-Version` header sent unconditionally would override every
  integration's credential-stamped version pin.
- A `429` that is thrown instead of retried turns a routine rate-limit trip
  into an outage in the caller's code.

None of those fail `go build`. All of them fail here.

HTTP is faked at `http.Client.Do` against the real default host
(`https://api.ro.am/v1`). URL construction, header merging, and retries all
run. An unhandled request fails the test rather than falling through to the
network.

| File | Pins |
| --- | --- |
| `pagination_test.go` | Cursor pagination walks two pages and stops: the `cursor` request field, the `nextCursor` response field, and the `groups` array are wired to the right names. |
| `version_header_test.go` | `Roam-Version` is sent when pinned at the client or per request, per-request wins, and it is **absent** when unset. Also pins `Token.Info(ctx, opts...)` — request options are not a second argument after an empty body. |
| `retry_test.go` | A `429` with `Retry-After: 2` is retried rather than thrown, and the wait is the header's value, not the client's backoff. Plus exhaustion, `WithoutRetries`, and that a `400` is not retried. |
| `webhook_verify_test.go` | Standard Webhooks verification as a consumer of `ro.am/roamhq/webhooks`: genuine `whsec_` deliveries accepted, tampered bodies rejected, multi-signature rotation, replay window. (The full unit suite lives in `webhooks/verify_test.go`.) |

This directory is listed in `.fernignore`. A regeneration must not delete it.
