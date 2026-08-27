package webhooks

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
)

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

func TestVerifySignatureAcceptsGenuineDelivery(t *testing.T) {
	secret, key := mintSecret(t)
	if err := VerifySignature([]byte(body), headersFor(key, messageID, nowUnix(), body), secret, nil); err != nil {
		t.Fatal(err)
	}
}

func TestVerifySignatureAcceptsSecretWithoutPrefix(t *testing.T) {
	secret, key := mintSecret(t)
	bare := strings.TrimPrefix(secret, "whsec_")
	if err := VerifySignature([]byte(body), headersFor(key, messageID, nowUnix(), body), bare, nil); err != nil {
		t.Fatal(err)
	}
}

func TestVerifySignatureRejectsTamperedBody(t *testing.T) {
	secret, key := mintSecret(t)
	err := VerifySignature([]byte(body+" "), headersFor(key, messageID, nowUnix(), body), secret, nil)
	if !isVerification(err) {
		t.Fatalf("got %v, want VerificationError", err)
	}
}

func TestVerifySignatureRejectsDifferentMessageID(t *testing.T) {
	secret, key := mintSecret(t)
	ts := nowUnix()
	h := make(http.Header)
	h.Set("webhook-id", messageID)
	h.Set("webhook-timestamp", ts)
	h.Set("webhook-signature", sign(key, "some-other-id", ts, body))
	err := VerifySignature([]byte(body), h, secret, nil)
	if !isVerification(err) {
		t.Fatalf("got %v, want VerificationError", err)
	}
}

func TestVerifySignatureRejectsWrongSecret(t *testing.T) {
	secret, _ := mintSecret(t)
	_, other := mintSecret(t)
	err := VerifySignature([]byte(body), headersFor(other, messageID, nowUnix(), body), secret, nil)
	if !isVerification(err) {
		t.Fatalf("got %v, want VerificationError", err)
	}
}

func TestVerifySignatureAcceptsSecondOfTwoSignatures(t *testing.T) {
	secret, key := mintSecret(t)
	_, stale := mintSecret(t)
	ts := nowUnix()
	h := make(http.Header)
	h.Set("webhook-id", messageID)
	h.Set("webhook-timestamp", ts)
	h.Set("webhook-signature", sign(stale, messageID, ts, body)+" "+sign(key, messageID, ts, body))
	if err := VerifySignature([]byte(body), h, secret, nil); err != nil {
		t.Fatal(err)
	}
}

func TestVerifySignatureIgnoresUnknownVersion(t *testing.T) {
	secret, key := mintSecret(t)
	ts := nowUnix()

	onlyUnknown := make(http.Header)
	onlyUnknown.Set("webhook-id", messageID)
	onlyUnknown.Set("webhook-timestamp", ts)
	onlyUnknown.Set("webhook-signature", "v2,"+sign(key, messageID, ts, body)[len("v1,"):])
	err := VerifySignature([]byte(body), onlyUnknown, secret, nil)
	if err == nil || !strings.Contains(err.Error(), "no v1 signature") {
		t.Fatalf("got %v, want no v1 signature", err)
	}

	mixed := make(http.Header)
	mixed.Set("webhook-id", messageID)
	mixed.Set("webhook-timestamp", ts)
	mixed.Set("webhook-signature", "v2,not-a-real-signature "+sign(key, messageID, ts, body))
	if err := VerifySignature([]byte(body), mixed, secret, nil); err != nil {
		t.Fatal(err)
	}
}

func TestVerifySignatureRejectsMissingHeader(t *testing.T) {
	secret, key := mintSecret(t)
	full := headersFor(key, messageID, nowUnix(), body)
	for _, omitted := range []string{"Webhook-Id", "Webhook-Timestamp", "Webhook-Signature"} {
		partial := full.Clone()
		partial.Del(omitted)
		err := VerifySignature([]byte(body), partial, secret, nil)
		if err == nil || !strings.Contains(strings.ToLower(err.Error()), "missing") {
			t.Fatalf("%s: got %v, want missing header", omitted, err)
		}
	}
}

func TestVerifySignatureReadsHeadersCaseInsensitively(t *testing.T) {
	secret, key := mintSecret(t)
	plain := headersFor(key, messageID, nowUnix(), body)
	upper := make(http.Header)
	upper["Webhook-Id"] = plain["Webhook-Id"]
	upper["Webhook-Timestamp"] = plain["Webhook-Timestamp"]
	upper["Webhook-Signature"] = plain["Webhook-Signature"]
	if err := VerifySignature([]byte(body), upper, secret, nil); err != nil {
		t.Fatal(err)
	}
}

func TestVerifySignatureRejectsEmptySecret(t *testing.T) {
	_, key := mintSecret(t)
	h := headersFor(key, messageID, nowUnix(), body)
	if err := VerifySignature([]byte(body), h, "", nil); err == nil || !strings.Contains(err.Error(), "secret is empty") {
		t.Fatalf("got %v, want secret is empty", err)
	}
	if err := VerifySignature([]byte(body), h, "whsec_", nil); err == nil || !strings.Contains(err.Error(), "did not base64-decode") {
		t.Fatalf("got %v, want did not base64-decode", err)
	}
}

func TestReplayWindow(t *testing.T) {
	secret, key := mintSecret(t)
	now := time.Now()
	ts := fmt.Sprintf("%d", now.Unix())
	h := headersFor(key, messageID, ts, body)

	err := VerifySignature([]byte(body), h, secret, &Options{
		Now: func() time.Time { return now.Add(400 * time.Second) },
	})
	if err == nil || !strings.Contains(err.Error(), "outside the 300s") {
		t.Fatalf("got %v, want outside the 300s", err)
	}

	if err := VerifySignature([]byte(body), h, secret, &Options{
		Now: func() time.Time { return now.Add(250 * time.Second) },
	}); err != nil {
		t.Fatal(err)
	}

	if err := VerifySignature([]byte(body), h, secret, &Options{
		Tolerance: 600 * time.Second,
		Now:       func() time.Time { return now.Add(400 * time.Second) },
	}); err != nil {
		t.Fatal(err)
	}

	bad := h.Clone()
	bad.Set("webhook-timestamp", "not-a-number")
	err = VerifySignature([]byte(body), bad, secret, nil)
	if err == nil || !strings.Contains(err.Error(), "is not a number") {
		t.Fatalf("got %v, want is not a number", err)
	}
}

func TestVerifyReturnsParsedBody(t *testing.T) {
	secret, key := mintSecret(t)
	event, err := Verify([]byte(body), headersFor(key, messageID, nowUnix(), body), secret, nil)
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

func TestVerifyRejectsNonJSON(t *testing.T) {
	secret, key := mintSecret(t)
	payload := "not json"
	_, err := Verify([]byte(payload), headersFor(key, messageID, nowUnix(), payload), secret, nil)
	if err == nil || !strings.Contains(err.Error(), "not JSON") {
		t.Fatalf("got %v, want not JSON", err)
	}
}

func isVerification(err error) bool {
	var v *VerificationError
	return errors.As(err, &v)
}
