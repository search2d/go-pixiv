package pixiv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type TokenProvider interface {
	Token(ctx context.Context) (string, error)
}

var DefaultOauthBaseURL = "https://oauth.secure.pixiv.net"

var DefaultOauthHeaders = map[string]string{
	"User-Agent":     "PixivAndroidApp/5.0.64 (Android 6.0; Google Nexus 5X - 6.0.0 - API 23 - 1080x1920)",
	"App-OS":         "android",
	"App-OS-Version": "6.0",
	"App-Version":    "5.0.64",
}

type OauthTokenProvider struct {
	Client     *http.Client
	BaseURL    string
	Headers    map[string]string
	Credential Credential
	Now        func() time.Time

	mx    sync.Mutex
	token *token
}

type Credential struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
}

func (p *OauthTokenProvider) Token(ctx context.Context) (string, error) {
	p.mx.Lock()
	defer p.mx.Unlock()

	if p.token == nil {
		if err := p.authorize(ctx); err != nil {
			return "", err
		}
		return p.token.accessToken, nil
	}

	if p.token.expired(p.now()) {
		if err := p.refresh(ctx); err != nil {
			return "", err
		}
		return p.token.accessToken, nil
	}

	return p.token.accessToken, nil
}

func (p *OauthTokenProvider) authorize(ctx context.Context) error {
	v := url.Values{}
	v.Set("username", p.Credential.Username)
	v.Set("password", p.Credential.Password)
	v.Set("client_id", p.Credential.ClientID)
	v.Set("client_secret", p.Credential.ClientSecret)
	v.Set("grant_type", "password")
	v.Set("get_secure_url", "true")

	req, err := http.NewRequest(http.MethodPost, p.baseURL()+"/auth/token", strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.request(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode <= 299) {
		return p.onFailure(res)
	}

	return p.onSuccess(res)
}

func (p *OauthTokenProvider) refresh(ctx context.Context) error {
	v := url.Values{}
	v.Set("refresh_token", p.token.refreshToken)
	v.Set("client_id", p.Credential.ClientID)
	v.Set("client_secret", p.Credential.ClientSecret)
	v.Set("grant_type", "refresh_token")
	v.Set("get_secure_url", "true")

	req, err := http.NewRequest(http.MethodPost, p.baseURL()+"/auth/token", strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.request(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode <= 299) {
		return p.onFailure(res)
	}

	return p.onSuccess(res)
}

func (p *OauthTokenProvider) request(req *http.Request) (*http.Response, error) {
	for k, v := range p.headers() {
		req.Header.Set(k, v)
	}

	return p.client().Do(req)
}

func (p *OauthTokenProvider) onSuccess(res *http.Response) error {
	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("Content-Type header = %q, should be \"application/json\"", res.Header.Get("Content-Type"))
	}

	var t Token

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return err
	}

	p.token = &token{
		accessToken:  t.Response.AccessToken,
		refreshToken: t.Response.RefreshToken,
		createdAt:    p.now(),
		expiresIn:    time.Duration(t.Response.ExpiresIn) * time.Second,
	}

	return nil
}

func (p *OauthTokenProvider) onFailure(res *http.Response) error {
	errToken := ErrToken{
		StatusCode: res.StatusCode,
		Status:     res.Status,
	}

	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		if err := json.NewDecoder(res.Body).Decode(&errToken.Body); err != nil {
			return err
		}
	}

	return errToken
}

func (p *OauthTokenProvider) client() *http.Client {
	if p.Client == nil {
		return http.DefaultClient
	}
	return p.Client
}

func (p *OauthTokenProvider) baseURL() string {
	if p.BaseURL == "" {
		return DefaultOauthBaseURL
	}
	return p.BaseURL
}

func (p *OauthTokenProvider) headers() map[string]string {
	if p.Headers == nil {
		return DefaultOauthHeaders
	}
	return p.Headers
}

func (p *OauthTokenProvider) now() time.Time {
	if p.Now == nil {
		return time.Now()
	}
	return p.Now()
}

type token struct {
	accessToken  string
	refreshToken string
	createdAt    time.Time
	expiresIn    time.Duration
}

func (t *token) expired(now time.Time) bool {
	expiredAt := t.createdAt.Add(t.expiresIn)
	return expiredAt.Equal(now) || expiredAt.Before(now)
}
