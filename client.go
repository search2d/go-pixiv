package pixiv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/search2d/go-pixiv/resp"
)

var defaultAPIBaseURL = "https://app-api.pixiv.net"

var defaultAPIHeaders = map[string]string{
	"User-Agent":     "PixivAndroidApp/5.0.64 (Android 6.0; Google Nexus 5X - 6.0.0 - API 23 - 1080x1920)",
	"App-OS":         "android",
	"App-OS-Version": "6.0",
	"App-Version":    "5.0.64",
}

type Client struct {
	client        *http.Client
	logger        *log.Logger
	tokenProvider TokenProvider
	baseURL       string
	headers       map[string]string
}

type ClientConfig struct {
	Client        *http.Client
	Logger        *log.Logger
	TokenProvider TokenProvider
	BaseURL       string
	Headers       map[string]string
}

func NewClient(cfg ClientConfig) *Client {
	c := &Client{tokenProvider: cfg.TokenProvider}

	if cfg.Client != nil {
		c.client = cfg.Client
	} else {
		c.client = http.DefaultClient
	}

	if cfg.Logger != nil {
		c.logger = cfg.Logger
	} else {
		c.logger = log.New(ioutil.Discard, "", 0)
	}

	if len(cfg.BaseURL) != 0 {
		c.baseURL = cfg.BaseURL
	} else {
		c.baseURL = defaultAPIBaseURL
	}

	if cfg.Headers != nil {
		c.headers = cfg.Headers
	} else {
		c.headers = defaultAPIHeaders
	}

	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	token, err := c.tokenProvider.Token()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return c.client.Do(req)
}

func (c *Client) onFailure(res *http.Response) error {
	r := ErrAPI{
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
	Date   *time.Time
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

func (p *GetIllustRankingParams) SetDate(date time.Time) *GetIllustRankingParams {
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

func (p *GetIllustRankingParams) validate() error {
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
		v.Set("date", p.Date.Format("2006-01-02"))
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

func (c *Client) GetIllustRanking(params *GetIllustRankingParams) (*resp.GetIllustRanking, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET",
		c.baseURL+"/v1/illust/ranking?"+params.buildQuery(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, c.onFailure(res)
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return nil, fmt.Errorf("Content-Type header = %q, should be \"application/json\"", res.Header.Get("Content-Type"))
	}

	var r resp.GetIllustRanking

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *Client) GetIllustRankingNext(nextURL string) (*resp.GetIllustRanking, error) {
	req, err := http.NewRequest("GET", nextURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, c.onFailure(res)
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return nil, fmt.Errorf("Content-Type header = %q, should be \"application/json\"", res.Header.Get("Content-Type"))
	}

	var r resp.GetIllustRanking

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
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

func (p *GetIllustDetailParams) validate() error {
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

func (c *Client) GetIllustDetail(params *GetIllustDetailParams) (*resp.GetIllustDetail, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET",
		c.baseURL+"/v1/illust/detail?"+params.buildQuery(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, c.onFailure(res)
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return nil, fmt.Errorf("Content-Type header = %q, should be \"application/json\"", res.Header.Get("Content-Type"))
	}

	var r resp.GetIllustDetail

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

type ErrAPI struct {
	StatusCode int
	Status     string
	JSON       string
}

func (e ErrAPI) Error() string {
	return fmt.Sprintf("%s", e.Status)
}

func (e ErrAPI) Decode() (resp.APIErrorBody, error) {
	var r resp.APIErrorBody

	if len(e.JSON) == 0 {
		return r, errors.New("no JSON contents")
	}

	if err := json.NewDecoder(strings.NewReader(e.JSON)).Decode(&r); err != nil {
		return r, err
	}

	return r, nil
}

type ErrInvalidParams struct {
	Errs []ErrInvalidParam
}

func (e *ErrInvalidParams) Error() string {
	return fmt.Sprintf("%d validation error(s) found", len(e.Errs))
}

func (e *ErrInvalidParams) Add(err ErrInvalidParam) {
	e.Errs = append(e.Errs, err)
}

func (e *ErrInvalidParams) Len() int {
	return len(e.Errs)
}

type ErrInvalidParam struct {
	Field   string
	Message string
}

func (e ErrInvalidParam) Error() string {
	return fmt.Sprintf("%s, %s", e.Field, e.Message)
}
