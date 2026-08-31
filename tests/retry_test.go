package contract

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	roamhq "ro.am/roamhq"
	"ro.am/roamhq/option"
)

// 429 handling and Retry-After.
//
// Pins the rate-limit modeling in fern/apis/roam/overlays.yml
// (developer-ro-am), which attaches a documented 429 response with a
// Retry-After header to every operation.
//
// The contract in docs/guides/sdks.md is specific: a 429 is retried, not
// thrown, and the wait comes from Retry-After rather than from the client's
// own backoff guess. Those are two separate claims, so this suite checks
// the delay value and not just the retry count — exponential backoff would
// also produce a second request.
//
// Timers are real. Fern's Go retrier applies no jitter on the Retry-After
// path, so Retry-After: 2 is a two-second sleep. Default exponential
// backoff for the first retry is ~1s, so 2s is what separates "honored the
// server" from "guessed and got lucky".

func always429(counter *int) func(*http.Request) *http.Response {
	return func(*http.Request) *http.Response {
		*counter++
		return jsonResponse(429, map[string]any{"ok": false, "error": "ratelimited"},
			http.Header{"Retry-After": []string{"2"}})
	}
}

func TestRetryAfterHonorsHeaderExactly(t *testing.T) {
	api := newFake(t)
	var attempts int
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		attempts++
		if attempts == 1 {
			return jsonResponse(429, map[string]any{"ok": false, "error": "ratelimited"},
				http.Header{"Retry-After": []string{"2"}})
		}
		return jsonResponse(200, json.RawMessage(`{
			"ok": true,
			"groups": [{"id": "g1", "name": "Engineering", "type": "standard"}],
			"nextCursor": null
		}`), nil)
	})

	start := time.Now()
	page, err := api.client().Group.List(context.Background(), &roamhq.ListGroupRequest{})
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
	if got := ids(page.Results); !equal(got, []string{"g1"}) {
		t.Fatalf("ids = %v", got)
	}
	if elapsed < 2*time.Second || elapsed >= 3*time.Second {
		t.Fatalf("elapsed = %s, want ~2s from Retry-After (not ~1s backoff)", elapsed)
	}
}

func TestRetryExhaustionThrowsTooManyRequests(t *testing.T) {
	api := newFake(t)
	var n int
	api.handle(http.MethodGet, "/v1/group.list", always429(&n))

	_, err := api.client().Group.List(context.Background(), &roamhq.ListGroupRequest{})
	if err == nil {
		t.Fatal("want TooManyRequestsError")
	}
	var tooMany *roamhq.TooManyRequestsError
	if !errors.As(err, &tooMany) {
		t.Fatalf("err = %T %v, want TooManyRequestsError", err, err)
	}
	// The generated retrier treats defaultRetryAttempts (2) as a cap on
	// total HTTP calls, not on retries: attempt 0 and 1 run, attempt 2
	// returns the previous error without a third request.
	if n != 2 {
		t.Fatalf("attempts = %d, want 2", n)
	}
}

func TestWithoutRetriesDoesNotRetry(t *testing.T) {
	api := newFake(t)
	var n int
	api.handle(http.MethodGet, "/v1/group.list", always429(&n))

	_, err := api.client().Group.List(
		context.Background(),
		&roamhq.ListGroupRequest{},
		option.WithoutRetries(),
	)
	if err == nil {
		t.Fatal("want error")
	}
	if n != 1 {
		t.Fatalf("attempts = %d, want 1", n)
	}
	// The generated retrier sleeps Retry-After *before* checking the
	// attempt cap, so WithoutRetries still waits, then returns. Pin the
	// request count; a Fern upgrade that skips the sleep is welcome and
	// will not fail this assertion.
}

func TestRetrySurfacesParsedErrorEnvelope(t *testing.T) {
	api := newFake(t)
	var n int
	api.handle(http.MethodGet, "/v1/group.list", always429(&n))

	_, err := api.client().Group.List(
		context.Background(),
		&roamhq.ListGroupRequest{},
		option.WithoutRetries(),
	)
	var tooMany *roamhq.TooManyRequestsError
	if !errors.As(err, &tooMany) {
		t.Fatalf("err = %T %v", err, err)
	}
	if tooMany.StatusCode != 429 {
		t.Fatalf("status = %d", tooMany.StatusCode)
	}
	if tooMany.Body == nil || tooMany.Body.Error != "ratelimited" {
		t.Fatalf("body = %+v, want error=ratelimited", tooMany.Body)
	}
}

func TestDoesNotRetry400(t *testing.T) {
	api := newFake(t)
	var n int
	api.handle(http.MethodGet, "/v1/group.list", func(*http.Request) *http.Response {
		n++
		return jsonResponse(400, map[string]any{"ok": false, "error": "invalid_arguments"}, nil)
	})

	_, err := api.client().Group.List(context.Background(), &roamhq.ListGroupRequest{})
	var bad *roamhq.BadRequestError
	if !errors.As(err, &bad) {
		t.Fatalf("err = %T %v, want BadRequestError", err, err)
	}
	if n != 1 {
		t.Fatalf("attempts = %d, want 1", n)
	}
}
