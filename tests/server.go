package contract

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"ro.am/roamhq/client"
	"ro.am/roamhq/option"
)

// The SDK's default environment. Tests intercept at http.Client.Do so
// URL construction, auth headers, and retries all run for real. A request
// that is not this host/path fails the test — same idea as msw's
// onUnhandledRequest: "error". Hitting api.ro.am from a contract test
// would be worse than no test.
const baseURL = "https://api.ro.am/v1"

const token = "rmk-test-token"

type seenReq struct {
	Method   string
	Path     string
	RawQuery string
	Header   http.Header
}

func (s seenReq) cursor() string {
	vals, _ := parseQuery(s.RawQuery)
	return vals.Get("cursor")
}

type fakeAPI struct {
	t        *testing.T
	mu       sync.Mutex
	seen     []seenReq
	handlers map[string]func(*http.Request) *http.Response
}

func newFake(t *testing.T) *fakeAPI {
	t.Helper()
	return &fakeAPI{
		t:        t,
		handlers: make(map[string]func(*http.Request) *http.Response),
	}
}

func (f *fakeAPI) handle(method, path string, h func(*http.Request) *http.Response) {
	f.handlers[method+" "+path] = h
}

func (f *fakeAPI) snapshot() []seenReq {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]seenReq, len(f.seen))
	copy(out, f.seen)
	return out
}

func (f *fakeAPI) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme != "https" || req.URL.Host != "api.ro.am" {
		f.t.Fatalf("unhandled request (wrong host): %s", req.URL)
	}
	key := req.Method + " " + req.URL.Path
	h, ok := f.handlers[key]
	if !ok {
		f.t.Fatalf("unhandled request: %s %s", req.Method, req.URL)
	}

	snap := seenReq{
		Method:   req.Method,
		Path:     req.URL.Path,
		RawQuery: req.URL.RawQuery,
		Header:   req.Header.Clone(),
	}
	f.mu.Lock()
	f.seen = append(f.seen, snap)
	f.mu.Unlock()

	return h(req), nil
}

func (f *fakeAPI) client(opts ...option.RequestOption) *client.Client {
	all := append([]option.RequestOption{
		option.WithToken(token),
		option.WithHTTPClient(f),
	}, opts...)
	return client.NewClient(all...)
}

func jsonResponse(status int, body any, hdr http.Header) *http.Response {
	raw, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	h := make(http.Header)
	for k, vs := range hdr {
		h[k] = vs
	}
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", "application/json")
	}
	return &http.Response{
		StatusCode:    status,
		Status:        http.StatusText(status),
		Header:        h,
		Body:          io.NopCloser(bytes.NewReader(raw)),
		ContentLength: int64(len(raw)),
	}
}

func cursorQuery(r *http.Request) string {
	return r.URL.Query().Get("cursor")
}

func parseQuery(raw string) (url.Values, error) {
	return url.ParseQuery(raw)
}
