package pixiv

import (
	"net/http"
	"testing"
)

func TestSetDownloadHeaders(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	SetDownloadHeaders(req)

	for k, v := range DefaultDownloadHeaders {
		if req.Header.Get(k) != v {
			t.Errorf("got %s header %q, want %q", k, req.Header.Get(k), v)
		}
	}
}
