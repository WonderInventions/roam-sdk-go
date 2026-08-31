package contract

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	roamhq "ro.am/roamhq"
	"ro.am/roamhq/core"
)

// Cursor pagination.
//
// Pins the x-fern-pagination modeling in fern/apis/roam/overrides.yml
// (developer-ro-am). The spec says `cursor` goes in on the request, the server
// hands back `nextCursor`, and the items live under a named array — for
// /group.list that array is `groups`. If any of those three names drift, the
// generated client silently stops paginating: it returns page one and reports
// no next page, which looks like "the account only has two groups" rather than
// like a bug.

func TestCursorPaginationWalksBothPagesAndStops(t *testing.T) {
	api := newFake(t)
	var seen []string
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		c := cursorQuery(r)
		seen = append(seen, c)
		switch c {
		case "":
			return jsonResponse(200, json.RawMessage(`{
				"ok": true,
				"groups": [
					{"id": "g1", "name": "Engineering", "type": "standard"},
					{"id": "g2", "name": "Design", "type": "standard"}
				],
				"nextCursor": "cursor-page-2"
			}`), nil)
		case "cursor-page-2":
			return jsonResponse(200, json.RawMessage(`{
				"ok": true,
				"groups": [
					{"id": "g3", "name": "Support", "type": "standard"}
				],
				"nextCursor": null
			}`), nil)
		default:
			t.Fatalf("unexpected cursor: %q", c)
			return nil
		}
	})

	c := api.client()
	page, err := c.Group.List(context.Background(), &roamhq.ListGroupRequest{})
	if err != nil {
		t.Fatal(err)
	}

	var names []string
	iter := page.Iterator()
	for iter.Next(context.Background()) {
		names = append(names, iter.Current().Name)
	}
	if err := iter.Err(); err != nil {
		t.Fatal(err)
	}

	if got, want := names, []string{"Engineering", "Design", "Support"}; !equal(got, want) {
		t.Fatalf("names = %v, want %v", got, want)
	}

	// Two requests, and the second one carried the cursor the first returned.
	// Asserting the cursor value — not just the request count — is what proves
	// next_cursor: $response.nextCursor is wired to the right field.
	if len(seen) != 2 || seen[0] != "" || seen[1] != "cursor-page-2" {
		t.Fatalf("cursors = %v, want ['', 'cursor-page-2']", seen)
	}
}

func TestCursorPaginationManualPageAPI(t *testing.T) {
	api := newFake(t)
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		switch cursorQuery(r) {
		case "":
			return jsonResponse(200, json.RawMessage(`{
				"ok": true,
				"groups": [
					{"id": "g1", "name": "Engineering", "type": "standard"},
					{"id": "g2", "name": "Design", "type": "standard"}
				],
				"nextCursor": "cursor-page-2"
			}`), nil)
		case "cursor-page-2":
			return jsonResponse(200, json.RawMessage(`{
				"ok": true,
				"groups": [
					{"id": "g3", "name": "Support", "type": "standard"}
				],
				"nextCursor": null
			}`), nil)
		default:
			t.Fatalf("unexpected cursor: %q", cursorQuery(r))
			return nil
		}
	})

	c := api.client()
	page, err := c.Group.List(context.Background(), &roamhq.ListGroupRequest{
		Limit: roamhq.Int(2),
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := ids(page.Results); !equal(got, []string{"g1", "g2"}) {
		t.Fatalf("page 1 ids = %v", got)
	}
	if page.RawResponse.Done {
		t.Fatal("page 1 should not be Done")
	}
	if page.Response.NextCursor == nil || *page.Response.NextCursor != "cursor-page-2" {
		t.Fatalf("page.response.nextCursor = %v, want cursor-page-2", page.Response.NextCursor)
	}

	page, err = page.GetNextPage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := ids(page.Results); !equal(got, []string{"g3"}) {
		t.Fatalf("page 2 ids = %v", got)
	}
	if !page.RawResponse.Done {
		t.Fatal("page 2 should be Done")
	}

	_, err = page.GetNextPage(context.Background())
	if !errors.Is(err, core.ErrNoPages) {
		t.Fatalf("GetNextPage after last page: %v, want ErrNoPages", err)
	}
}

func TestCursorPaginationDoesNotRequestAThirdPage(t *testing.T) {
	api := newFake(t)
	var n int
	api.handle(http.MethodGet, "/v1/group.list", func(r *http.Request) *http.Response {
		n++
		if cursorQuery(r) == "" {
			return jsonResponse(200, json.RawMessage(`{
				"ok": true,
				"groups": [{"id": "g1", "name": "Engineering", "type": "standard"}],
				"nextCursor": "cursor-page-2"
			}`), nil)
		}
		return jsonResponse(200, json.RawMessage(`{
			"ok": true,
			"groups": [{"id": "g3", "name": "Support", "type": "standard"}],
			"nextCursor": null
		}`), nil)
	})

	c := api.client()
	page, err := c.Group.List(context.Background(), &roamhq.ListGroupRequest{})
	if err != nil {
		t.Fatal(err)
	}
	iter := page.Iterator()
	for iter.Next(context.Background()) {
		_ = iter.Current()
	}
	if err := iter.Err(); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("requests = %d, want 2", n)
	}
}

func ids(groups []*roamhq.ListGroupResponseGroupsItem) []string {
	out := make([]string, len(groups))
	for i, g := range groups {
		out[i] = g.ID
	}
	return out
}

func equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
