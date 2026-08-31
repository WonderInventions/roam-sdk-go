package contract

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	roamhq "ro.am/roamhq"
	"ro.am/roamhq/option"
)

// The Roam-Version request header.
//
// Pins x-fern-global-headers in fern/apis/roam/overrides.yml and the matching
// headers: block in fern/apis/roam/generators.yml (developer-ro-am).
//
// Two halves matter, and the second is the easy one to lose. Sending the
// header when the caller asks for a pin is obvious. Not sending it otherwise
// is the subtle part: Roam falls back to the version stamped on the
// credential when the header is absent, so a client that always sent
// something — an empty string, or a baked-in default — would silently
// override every integration's pin.

const (
	pinned = "2026-07-23"
	other  = "2026-01-15"
)

func emptyGroupList() json.RawMessage {
	return json.RawMessage(`{"ok": true, "groups": [], "nextCursor": null}`)
}

func TestRoamVersionAbsentWhenNotPinned(t *testing.T) {
	api := newFake(t)
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		return jsonResponse(200, emptyGroupList(), nil)
	})

	if _, err := api.client().Group.List(context.Background(), &roamhq.ListGroupRequest{}); err != nil {
		t.Fatal(err)
	}
	if api.snapshot()[0].Header.Get("Roam-Version") != "" {
		t.Fatalf("Roam-Version = %q, want absent", api.snapshot()[0].Header.Get("Roam-Version"))
	}
}

func TestRoamVersionSentWhenSetOnClient(t *testing.T) {
	api := newFake(t)
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		return jsonResponse(200, emptyGroupList(), nil)
	})

	c := api.client(option.WithRoamVersion(roamhq.String(pinned)))
	if _, err := c.Group.List(context.Background(), &roamhq.ListGroupRequest{}); err != nil {
		t.Fatal(err)
	}
	if got := api.snapshot()[0].Header.Get("Roam-Version"); got != pinned {
		t.Fatalf("Roam-Version = %q, want %q", got, pinned)
	}
}

func TestRoamVersionSentWhenSetOnRequest(t *testing.T) {
	api := newFake(t)
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		return jsonResponse(200, emptyGroupList(), nil)
	})

	if _, err := api.client().Group.List(
		context.Background(),
		&roamhq.ListGroupRequest{},
		option.WithRoamVersion(roamhq.String(pinned)),
	); err != nil {
		t.Fatal(err)
	}
	if got := api.snapshot()[0].Header.Get("Roam-Version"); got != pinned {
		t.Fatalf("Roam-Version = %q, want %q", got, pinned)
	}
}

func TestRoamVersionPerRequestOverridesClient(t *testing.T) {
	api := newFake(t)
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		return jsonResponse(200, emptyGroupList(), nil)
	})

	c := api.client(option.WithRoamVersion(roamhq.String(pinned)))
	if _, err := c.Group.List(
		context.Background(),
		&roamhq.ListGroupRequest{},
		option.WithRoamVersion(roamhq.String(other)),
	); err != nil {
		t.Fatal(err)
	}
	if got := api.snapshot()[0].Header.Get("Roam-Version"); got != other {
		t.Fatalf("Roam-Version = %q, want %q", got, other)
	}
}

func TestRoamVersionClientPinKeptWhenNotOverridden(t *testing.T) {
	api := newFake(t)
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		return jsonResponse(200, emptyGroupList(), nil)
	})

	c := api.client(option.WithRoamVersion(roamhq.String(pinned)))
	if _, err := c.Group.List(context.Background(), &roamhq.ListGroupRequest{}, option.WithRoamVersion(roamhq.String(other))); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Group.List(context.Background(), &roamhq.ListGroupRequest{}); err != nil {
		t.Fatal(err)
	}
	seen := api.snapshot()
	if got := seen[0].Header.Get("Roam-Version"); got != other {
		t.Fatalf("first Roam-Version = %q, want %q", got, other)
	}
	if got := seen[1].Header.Get("Roam-Version"); got != pinned {
		t.Fatalf("second Roam-Version = %q, want %q", got, pinned)
	}
}

func TestRoamVersionStillSendsBearerAuth(t *testing.T) {
	api := newFake(t)
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		return jsonResponse(200, emptyGroupList(), nil)
	})

	c := api.client(option.WithRoamVersion(roamhq.String(pinned)))
	if _, err := c.Group.List(context.Background(), &roamhq.ListGroupRequest{}); err != nil {
		t.Fatal(err)
	}
	if got := api.snapshot()[0].Header.Get("Authorization"); got != "Bearer "+token {
		t.Fatalf("Authorization = %q", got)
	}
}

func TestTokenInfoTakesRequestOptionsAsFirstArgumentAfterCtx(t *testing.T) {
	// token.info takes no request payload, so the generated signature is
	// Info(ctx, opts...) — request options are not a second argument after
	// an empty body. Worth pinning: it is an easy mistake to write, it
	// fails quietly rather than loudly, and the shape differs from every
	// method that does take a body.
	api := newFake(t)
	api.handle(http.MethodGet, "/v1/token.info", func(r *http.Request) *http.Response {
		return jsonResponse(200, json.RawMessage(`{"ok": true, "tokenType": "org"}`), nil)
	})

	if _, err := api.client().Token.Info(context.Background(), option.WithRoamVersion(roamhq.String(pinned))); err != nil {
		t.Fatal(err)
	}
	if got := api.snapshot()[0].Header.Get("Roam-Version"); got != pinned {
		t.Fatalf("Roam-Version = %q, want %q", got, pinned)
	}
}
