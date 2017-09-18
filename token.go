package pixiv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/search2d/go-pixiv/resp"
)

type TokenProvider interface {
	Token() (string, error)
}

var defaultOauthBaseURL = "https://oauth.secure.pixiv.net"

var defaultOauthHeaders = map[string]string{
	"User-Agent":     "PixivAndroidApp/5.0.64 (Android 6.0; Google Nexus 5X - 6.0.0 - API 23 - 1080x1920)",
	"App-OS":         "android",
	"App-OS-Version": "6.0",
	"App-Version":    "5.0.64",
}

var now = func() time.Time {
	return time.Now()
}

type OauthTokenProvider struct {
	client     *http.Client
	logger     *log.Logger
	baseURL    string
	headers    map[string]string
	credential Credential

	mx    sync.Mutex
	token *token
}

type OauthTokenProviderConfig struct {
	Client     *http.Client
	Logger     *log.Logger
	BaseURL    string
	Headers    map[string]string
	Credential Credential
}

type Credential struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
}

func NewOauthTokenProvider(cfg OauthTokenProviderConfig) *OauthTokenProvider {
	p := &OauthTokenProvider{credential: cfg.Credential}

	if cfg.Client != nil {
		p.client = cfg.Client
	} else {
		p.client = http.DefaultClient
	}

	if cfg.Logger != nil {
		p.logger = cfg.Logger
	} else {
		p.logger = log.New(ioutil.Discard, "", 0)
	}

	if len(cfg.BaseURL) != 0 {
		p.baseURL = cfg.BaseURL
	} else {
		p.baseURL = defaultOauthBaseURL
	}

	if cfg.Headers != nil {
		p.headers = cfg.Headers
	} else {
		p.headers = defaultOauthHeaders
	}

	return p
}

func (p *OauthTokenProvider) Token() (string, error) {
	p.mx.Lock()
	defer p.mx.Unlock()

	if p.token == nil {
		if err := p.authorize(); err != nil {
			return "", err
		}
		return p.token.accessToken, nil
	}

	if p.token.expired() {
		if err := p.refresh(); err != nil {
			return "", err
		}
		return p.token.accessToken, nil
	}

	return p.token.accessToken, nil
}

func (p *OauthTokenProvider) authorize() error {
	v := url.Values{}
	v.Set("username", p.credential.Username)
	v.Set("password", p.credential.Password)
	v.Set("client_id", p.credential.ClientID)
	v.Set("client_secret", p.credential.ClientSecret)
	v.Set("grant_type", "password")
	v.Set("get_secure_url", "true")

	req, err := http.NewRequest(
		"POST",
		p.baseURL+"/auth/token",
		strings.NewReader(v.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.request(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode <= 299) {
		return p.onFailure(res)
	}

	return p.onSuccess(res)
}

func (p *OauthTokenProvider) refresh() error {
	v := url.Values{}
	v.Set("refresh_token", p.token.refreshToken)
	v.Set("client_id", p.credential.ClientID)
	v.Set("client_secret", p.credential.ClientSecret)
	v.Set("grant_type", "refresh_token")
	v.Set("get_secure_url", "true")

	req, err := http.NewRequest(
		"POST",
		p.baseURL+"/auth/token",
		strings.NewReader(v.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.request(req)
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
	for k, v := range p.headers {
		req.Header.Set(k, v)
	}

	return p.client.Do(req)
}

func (p *OauthTokenProvider) onSuccess(res *http.Response) error {
	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("Content-Type header = %q, should be \"application/json\"", res.Header.Get("Content-Type"))
	}

	var r resp.Token

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}

	p.token = &token{
		accessToken:  r.Response.AccessToken,
		refreshToken: r.Response.RefreshToken,
		createdOn:    now(),
		expiresIn:    time.Duration(r.Response.ExpiresIn) * time.Second,
	}

	return nil
}

func (p *OauthTokenProvider) onFailure(res *http.Response) error {
	r := ErrToken{
		StatusCode: res.StatusCode,
		Status:     res.Status,
	}

	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		r.JSON = string(buf)
	}

	return r
}

type token struct {
	accessToken  string
	refreshToken string
	createdOn    time.Time
	expiresIn    time.Duration
}

func (t *token) expired() bool {
	return t.createdOn.Add(t.expiresIn).After(now())
}

type ErrToken struct {
	StatusCode int
	Status     string
	JSON       string
}

func (e ErrToken) Error() string {
	return fmt.Sprintf("%s", e.Status)
}

func (e ErrToken) TokenErrorBody() (resp.TokenErrorBody, error) {
	var r resp.TokenErrorBody

	if len(e.JSON) == 0 {
		return r, errors.New("no JSON contents")
	}

	if err := json.NewDecoder(strings.NewReader(e.JSON)).Decode(&r); err != nil {
		return r, err
	}

	return r, nil
}
