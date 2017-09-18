package pixiv

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/search2d/go-pixiv/resp"
)

func TestOauthTokenProvider_Token(t *testing.T) {
	cnt := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			cnt++
		}()

		for k, v := range defaultOauthHeaders {
			if r.Header.Get(k) != v {
				t.Errorf("got %s header = %q, want %q", k, r.Header.Get(k), v)
			}
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}

		switch cnt {
		case 0:
			if g, e := r.URL.Path, "/auth/token"; g != e {
				t.Errorf("got URL path %q, want %q", g, e)
			}

			if g, e := r.Method, http.MethodPost; g != e {
				t.Errorf("got HTTP method %q, want %q", g, e)
			}

			if g, e := r.Header.Get("Content-Type"), "application/x-www-form-urlencoded"; g != e {
				t.Errorf("got Content-Type header = %q, want %q", g, e)
			}

			expectedForm := url.Values{
				"username":       []string{"USERNAME"},
				"password":       []string{"PASSWORD"},
				"client_id":      []string{"CLIENT_ID"},
				"client_secret":  []string{"CLIENT_SECRET"},
				"grant_type":     []string{"password"},
				"get_secure_url": []string{"true"},
			}
			if g, e := r.Form, expectedForm; !reflect.DeepEqual(g, e) {
				t.Errorf("got form %#v, want %#v", g, e)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(fixture("fixtures/token_authorize.json"))
		case 1:
			if g, e := r.URL.Path, "/auth/token"; g != e {
				t.Errorf("got URL path %q, want %q", g, e)
			}

			if g, e := r.Method, http.MethodPost; g != e {
				t.Errorf("got HTTP method %q, want %q", g, e)
			}

			if g, e := r.Header.Get("Content-Type"), "application/x-www-form-urlencoded"; g != e {
				t.Errorf("got Content-Type header = %q, want %q", g, e)
			}

			expectedForm := url.Values{
				"refresh_token":  []string{"wgNv1gZ0y8Z1nIyG4bRbpT2yNMs3hvHhHLIhXDc47G8"},
				"client_id":      []string{"CLIENT_ID"},
				"client_secret":  []string{"CLIENT_SECRET"},
				"grant_type":     []string{"refresh_token"},
				"get_secure_url": []string{"true"},
			}
			if g, e := r.Form, expectedForm; !reflect.DeepEqual(g, e) {
				t.Errorf("got form %#v, want %#v", g, e)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(fixture("fixtures/token_refresh.json"))
		default:
			t.Fatal("too many requests")
		}
	}))
	defer ts.Close()

	tp := NewOauthTokenProvider(OauthTokenProviderConfig{
		BaseURL: ts.URL,
		Credential: Credential{
			Username:     "USERNAME",
			Password:     "PASSWORD",
			ClientID:     "CLIENT_ID",
			ClientSecret: "CLIENT_SECRET",
		},
	})

	now = func() time.Time {
		return time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	token1, err := tp.Token()
	if err != nil {
		t.Fatal(err)
	}

	if g, e := token1, "ATN7bmWC7Kg1OneEqSPa9GxKm1l1uVHa8cQQKme7BGY"; g != e {
		t.Errorf("got token %q, want %q", g, e)
	}

	now = func() time.Time {
		return time.Date(2017, 1, 1, 0, 30, 0, 0, time.UTC)
	}

	token2, err := tp.Token()
	if err != nil {
		t.Fatal(err)
	}

	if g, e := token2, "cIPvPp368gKDU4DP7sXhbFzqKiXrGpwFJrbXF40fpUY"; g != e {
		t.Errorf("got token %q, want %q", g, e)
	}
}

func TestOauthTokenProvider_Token_BadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(fixture("fixtures/token_error.json"))
	}))
	defer ts.Close()

	tp := NewOauthTokenProvider(OauthTokenProviderConfig{
		BaseURL: ts.URL,
		Credential: Credential{
			Username:     "USERNAME",
			Password:     "PASSWORD",
			ClientID:     "CLIENT_ID",
			ClientSecret: "CLIENT_SECRET",
		},
	})

	_, err := tp.Token()

	if err == nil {
		t.Fatalf("Token() should return an error if 400 response is received")
	}

	errToken, ok := err.(ErrToken)
	if !ok {
		t.Fatalf("Token() should return an ErrToken if 404 response is received")
	}

	if g, e := errToken.StatusCode, http.StatusBadRequest; g != e {
		t.Errorf("got StatusCode %v, want %v", g, e)
	}

	tokenErrorBody, err := errToken.TokenErrorBody()
	if err != nil {
		t.Fatal(err)
	}

	expectedTokenErrorBody := resp.TokenErrorBody{
		HasError: true,
		Errors: map[string]resp.TokenError{
			"system": {
				Message: "103:pixiv ID、またはメールアドレス、パスワードが正しいかチェックしてください。",
				Code:    1508,
			},
		},
	}
	if g, e := tokenErrorBody, expectedTokenErrorBody; !reflect.DeepEqual(g, e) {
		t.Errorf("got TokenErrorBody %#v, want %#v", g, e)
	}

}
