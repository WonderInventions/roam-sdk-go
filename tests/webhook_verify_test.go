package contract

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"ro.am/roamhq/webhooks"
)

// Standard Webhooks signature verification, tested as a consumer of
// ro.am/roamhq/webhooks — the same import path integrators use.
//
// The first test is the one that matters most: it uses a secret shaped
// exactly the way Roam mints them, with a random 32-byte key. That is the
// case Fern's generated helper fails, because it keys the HMAC with the
// UTF-8 bytes of the whsec_ string instead of the decoded bytes.

const (
	messageID = "b0d9f6d2-1c3c-4a7e-9c1b-2d5e8f0a1b2c"
	body      = `{"eventId":"b0d9f6d2-1c3c-4a7e-9c1b-2d5e8f0a1b2c","type":"chat.message"}`
)

func mintSecret(t *testing.T) (secret string, key []byte) {
	t.Helper()
	key = make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	return "whsec_" + base64.StdEncoding.EncodeToString(key), key
}

func sign(key []byte, id, ts, payload string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(id + "." + ts + "." + payload))
	return "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func nowUnix() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func headersFor(key []byte, id, ts, payload string) http.Header {
	h := make(http.Header)
	h.Set("webhook-id", id)
	h.Set("webhook-timestamp", ts)
	h.Set("webhook-signature", sign(key, id, ts, payload))
	return h
}

func TestVerifyAcceptsGenuineWhsecDelivery(t *testing.T) {
	secret, key := mintSecret(t)
	if err := webhooks.VerifySignature([]byte(body), headersFor(key, messageID, nowUnix(), body), secret, nil); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyAcceptsSecretWithoutPrefix(t *testing.T) {
	secret, key := mintSecret(t)
	bare := strings.TrimPrefix(secret, "whsec_")
	if err := webhooks.VerifySignature([]byte(body), headersFor(key, messageID, nowUnix(), body), bare, nil); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyRejectsTamperedBody(t *testing.T) {
	secret, key := mintSecret(t)
	err := webhooks.VerifySignature([]byte(body+" "), headersFor(key, messageID, nowUnix(), body), secret, nil)
	var v *webhooks.VerificationError
	if !errors.As(err, &v) {
		t.Fatalf("got %v, want VerificationError", err)
	}
}

func TestVerifyAcceptsSecondOfTwoSignatures(t *testing.T) {
	secret, key := mintSecret(t)
	_, stale := mintSecret(t)
	ts := nowUnix()
	h := make(http.Header)
	h.Set("webhook-id", messageID)
	h.Set("webhook-timestamp", ts)
	h.Set("webhook-signature", sign(stale, messageID, ts, body)+" "+sign(key, messageID, ts, body))
	if err := webhooks.VerifySignature([]byte(body), h, secret, nil); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyReplayWindow(t *testing.T) {
	secret, key := mintSecret(t)
	now := time.Now()
	ts := fmt.Sprintf("%d", now.Unix())
	h := headersFor(key, messageID, ts, body)

	err := webhooks.VerifySignature([]byte(body), h, secret, &webhooks.Options{
		Now: func() time.Time { return now.Add(400 * time.Second) },
	})
	if err == nil || !strings.Contains(err.Error(), "outside the 300s") {
		t.Fatalf("got %v, want outside the 300s", err)
	}

	if err := webhooks.VerifySignature([]byte(body), h, secret, &webhooks.Options{
		Now: func() time.Time { return now.Add(250 * time.Second) },
	}); err != nil {
		t.Fatal(err)
	}

	if err := webhooks.VerifySignature([]byte(body), h, secret, &webhooks.Options{
		Tolerance: 600 * time.Second,
		Now:       func() time.Time { return now.Add(400 * time.Second) },
	}); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyReturnsParsedBody(t *testing.T) {
	secret, key := mintSecret(t)
	event, err := webhooks.Verify([]byte(body), headersFor(key, messageID, nowUnix(), body), secret, nil)
	if err != nil {
		t.Fatal(err)
	}
	var parsed struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(event, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed.Type != "chat.message" {
		t.Fatalf("type = %q", parsed.Type)
	}
}
