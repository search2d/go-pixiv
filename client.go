package pixiv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var DefaultAPIBaseURL = "https://app-api.pixiv.net"

var DefaultAPIHeaders = map[string]string{
	"User-Agent":     "PixivAndroidApp/5.0.64 (Android 6.0; Google Nexus 5X - 6.0.0 - API 23 - 1080x1920)",
	"App-OS":         "android",
	"App-OS-Version": "6.0",
	"App-Version":    "5.0.64",
}

type Client struct {
	Client        *http.Client
	BaseURL       string
	Headers       map[string]string
	TokenProvider TokenProvider
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	token, err := c.TokenProvider.Token(req.Context())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	for k, v := range c.headers() {
		req.Header.Set(k, v)
	}

	return c.client().Do(req)
}

func (c *Client) client() *http.Client {
	if c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

func (c *Client) baseURL() string {
	if c.BaseURL == "" {
		return DefaultAPIBaseURL
	}
	return c.BaseURL
}

func (c *Client) headers() map[string]string {
	if c.Headers == nil {
		return DefaultAPIHeaders
	}
	return c.Headers
}

func (c *Client) onSuccess(res *http.Response, val interface{}) error {
	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("Content-Type header = %q, should be \"application/json\"", res.Header.Get("Content-Type"))
	}

	return json.NewDecoder(res.Body).Decode(val)
}

func (c *Client) onFailure(res *http.Response) error {
	errAPI := ErrAPI{StatusCode: res.StatusCode, Status: res.Status}

	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		if err := json.NewDecoder(res.Body).Decode(&errAPI.Body); err != nil {
			return err
		}
	}

	return errAPI
}

const (
	RankingModeDay          = "day"
	RankingModeDayMale      = "day_male"
	RankingModeDayFemale    = "day_female"
	RankingModeDayR18       = "day_r18"
	RankingModeDayMaleR18   = "day_male_r18"
	RankingModeDayFemaleR18 = "day_female_r18"
	RankingModeWeek         = "week"
	RankingModeWeekOriginal = "week_original"
	RankingModeWeekRookie   = "week_rookie"
	RankingModeWeekR18      = "week_r18"
	RankingModeWeekR18G     = "week_r18g"
	RankingModeMonth        = "month"
)

type GetIllustRankingParams struct {
	Mode   *string
	Date   *string
	Offset *int
	Filter *string
}

func NewGetIllustRankingParams() *GetIllustRankingParams {
	return &GetIllustRankingParams{}
}

func (p *GetIllustRankingParams) SetMode(mode string) *GetIllustRankingParams {
	p.Mode = &mode
	return p
}

func (p *GetIllustRankingParams) SetDate(date string) *GetIllustRankingParams {
	p.Date = &date
	return p
}

func (p *GetIllustRankingParams) SetOffset(offset int) *GetIllustRankingParams {
	p.Offset = &offset
	return p
}

func (p *GetIllustRankingParams) SetFilter(filter string) *GetIllustRankingParams {
	p.Filter = &filter
	return p
}

func (p *GetIllustRankingParams) Validate() error {
	err := &ErrInvalidParams{}

	if p.Mode == nil {
		err.Add(ErrInvalidParam{"Mode", "missing required field"})
	}

	if err.Len() > 0 {
		return err
	}

	return nil
}

func (p *GetIllustRankingParams) buildQuery() string {
	v := url.Values{}

	v.Set("mode", *p.Mode)

	if p.Date != nil {
		v.Set("date", *p.Date)
	}

	if p.Offset != nil {
		v.Set("offset", strconv.Itoa(*p.Offset))
	}

	if p.Filter != nil {
		v.Set("filter", *p.Filter)
	} else {
		v.Set("filter", "for_android")
	}

	return v.Encode()
}

func (c *Client) GetIllustRanking(ctx context.Context, params *GetIllustRankingParams) (*GetIllustRanking, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodGet,
		c.baseURL()+"/v1/illust/ranking?"+params.buildQuery(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, c.onFailure(res)
	}

	var result GetIllustRanking

	if err := c.onSuccess(res, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetIllustRankingNext(ctx context.Context, nextURL string) (*GetIllustRanking, error) {
	req, err := http.NewRequest(http.MethodGet, nextURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, c.onFailure(res)
	}

	var ranking GetIllustRanking

	if err := c.onSuccess(res, &ranking); err != nil {
		return nil, err
	}

	return &ranking, nil
}

type GetIllustDetailParams struct {
	IllustID *int
}

func NewGetIllustDetailParams() *GetIllustDetailParams {
	return &GetIllustDetailParams{}
}

func (p *GetIllustDetailParams) SetIllustID(illustID int) *GetIllustDetailParams {
	p.IllustID = &illustID
	return p
}

func (p *GetIllustDetailParams) Validate() error {
	err := &ErrInvalidParams{}

	if p.IllustID == nil {
		err.Add(ErrInvalidParam{"IllustID", "missing required field"})
	}

	if err.Len() > 0 {
		return err
	}

	return nil
}

func (p *GetIllustDetailParams) buildQuery() string {
	v := url.Values{}

	v.Set("illust_id", strconv.Itoa(*p.IllustID))

	return v.Encode()
}

func (c *Client) GetIllustDetail(ctx context.Context, params *GetIllustDetailParams) (*GetIllustDetail, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodGet,
		c.baseURL()+"/v1/illust/detail?"+params.buildQuery(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, c.onFailure(res)
	}

	var result GetIllustDetail

	if err := c.onSuccess(res, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
