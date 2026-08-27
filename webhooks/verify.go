// Package webhooks verifies Roam webhook deliveries.
//
// HAND-WRITTEN. Not generated, and listed in .fernignore — the
// "Regenerate SDK from OpenAPI spec" workflow replaces the rest of this
// module wholesale with `rsync --delete`, so security-relevant code must
// not live in a generated package.
//
// Why this exists rather than the generated helper: Fern's webhook signature
// model has no way to describe how a signing secret is encoded. Its hmac
// config covers the header, algorithm, output encoding, payload format, and
// timestamp — but the secret itself is always passed to the HMAC as a raw
// string. Roam (like every Standard Webhooks implementation) issues secrets as
// `whsec_<base64>` and keys the HMAC with the decoded bytes. The generated
// helper therefore rejects every genuine delivery.
//
// Verified against vectors computed the way the Roam server computes them.
//
//	signed content = "{webhook-id}.{webhook-timestamp}.{raw body}"
//	HMAC-SHA256, base64, keyed with base64_decode(secret minus "whsec_")
//	header value = "v1,<signature>", space-separated for multiple signatures
//
// See https://www.standardwebhooks.com/
package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTolerance = 300 * time.Second
	secretPrefix     = "whsec_"
	supportedVersion = "v1"
)

// VerificationError is returned when a delivery cannot be verified.
// The message says which check failed. Treat any error from this package
// as a 401 — never fall through to processing the payload.
type VerificationError struct {
	msg string
}

func (e *VerificationError) Error() string { return e.msg }

func fail(format string, args ...any) error {
	return &VerificationError{msg: fmt.Sprintf(format, args...)}
}

// Options tune verification. The zero value is the Standard Webhooks default.
type Options struct {
	// Tolerance is how far the delivery's timestamp may be from now.
	// Defaults to 300s, the Standard Webhooks default and what Roam's own
	// receivers use.
	//
	// Roam signs once and reuses the same signature across delivery retries,
	// and the retry ladder runs to roughly six minutes — so a final retry can
	// arrive outside a strict 300s window. Widen this if you would rather
	// accept a late retry than drop it.
	Tolerance time.Duration

	// Now overrides the clock. Tests use this; callers should not.
	Now func() time.Time
}

// Verify a webhook delivery and return its parsed JSON body.
//
// payload must be the raw request body, exactly as received. Parsing and
// re-serializing changes the bytes (key order, whitespace, unicode escapes)
// and the signature will not match.
//
//	body, err := io.ReadAll(r.Body)
//	event, err := webhooks.Verify(body, r.Header, os.Getenv("ROAM_WEBHOOK_SECRET"))
func Verify(payload []byte, headers http.Header, secret string, opts *Options) (json.RawMessage, error) {
	if err := VerifySignature(payload, headers, secret, opts); err != nil {
		return nil, err
	}
	if !json.Valid(payload) {
		return nil, fail("signature is valid but the body is not JSON")
	}
	return json.RawMessage(append([]byte(nil), payload...)), nil
}

// VerifySignature checks a delivery without parsing the body.
//
// Use this when the payload is not JSON, or when you want to hand the raw
// bytes to your own parser.
func VerifySignature(payload []byte, headers http.Header, secret string, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}
	messageID := headers.Get("webhook-id")
	if messageID == "" {
		return fail("missing webhook-id header")
	}
	timestamp := headers.Get("webhook-timestamp")
	if timestamp == "" {
		return fail("missing webhook-timestamp header")
	}
	signatureHeader := headers.Get("webhook-signature")
	if signatureHeader == "" {
		return fail("missing webhook-signature header")
	}

	tolerance := opts.Tolerance
	if tolerance == 0 {
		tolerance = defaultTolerance
	}
	now := time.Now
	if opts.Now != nil {
		now = opts.Now
	}
	if err := assertTimestampIsFresh(timestamp, tolerance, now()); err != nil {
		return err
	}

	key, err := decodeSecret(secret)
	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(messageID))
	mac.Write([]byte("."))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(payload)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// The header carries a space-separated list of versioned signatures, so
	// that a secret can be rotated without dropping deliveries: during
	// rotation the sender signs with both the old and new secret. Accept the
	// delivery if ANY v1 entry matches. Ignoring the extras — as a
	// single-signature reader does — silently breaks every rotation.
	candidates := 0
	matched := false
	for _, part := range strings.Fields(signatureHeader) {
		version, sig, ok := strings.Cut(part, ",")
		if !ok {
			version, sig = supportedVersion, part
		}
		if version != supportedVersion {
			continue
		}
		candidates++
		// Compare every candidate rather than stopping at the first match, so
		// the work done does not depend on which signature matched.
		if hmac.Equal([]byte(sig), []byte(expected)) {
			matched = true
		}
	}
	if candidates == 0 {
		return fail("no %s signature found in the webhook-signature header", supportedVersion)
	}
	if !matched {
		return fail("webhook signature does not match the computed signature")
	}
	return nil
}

// Roam issues signing secrets as `whsec_<base64>`, and the HMAC key is the
// decoded bytes — not the string. This is the step the generated helper omits,
// and the entire reason this package exists.
//
// The prefix is optional so that a secret copied without it still works.
func decodeSecret(secret string) ([]byte, error) {
	if secret == "" {
		return nil, fail("signing secret is empty")
	}
	encoded := strings.TrimPrefix(secret, secretPrefix)
	key, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(key) == 0 {
		return nil, fail("signing secret did not base64-decode to any bytes; expected the whsec_… value from Roam Administration → Developer")
	}
	return key, nil
}

func assertTimestampIsFresh(timestamp string, tolerance time.Duration, now time.Time) error {
	sent, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fail("webhook-timestamp is not a number: %s", timestamp)
	}
	skew := now.Sub(time.Unix(sent, 0))
	if skew < 0 {
		skew = -skew
	}
	if skew > tolerance {
		return fail("webhook-timestamp is %ss away from now, outside the %ss tolerance",
			formatSeconds(skew), formatSeconds(tolerance))
	}
	return nil
}

func formatSeconds(d time.Duration) string {
	return strconv.FormatInt(int64(math.Round(d.Seconds())), 10)
}
