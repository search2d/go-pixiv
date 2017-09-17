package pixiv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/search2d/go-pixiv/resp"
)

type mockTokenProvider struct {
	token string
	err   error
}

func (p *mockTokenProvider) Token() (string, error) {
	return p.token, p.err
}

func TestClient_GetIllustRanking(t *testing.T) {
	tp := &mockTokenProvider{token: "ATN7bmWC7Kg1OneEqSPa9GxKm1l1uVHa8cQQKme7BGY"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g, e := r.URL.Path, "/v1/illust/ranking"; g != e {
			t.Errorf("got URL path %q, want %q", g, e)
		}

		if g, e := r.Method, http.MethodGet; g != e {
			t.Errorf("got HTTP method %q, want %q", g, e)
		}

		if g, e := r.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", tp.token); g != e {
			t.Errorf("got Authorization header = %q, want %q", g, e)
		}

		for k, v := range defaultAPIHeaders {
			if r.Header.Get(k) != v {
				t.Errorf("got %s header = %q, want %q", k, r.Header.Get(k), v)
			}
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}

		expectedForm := url.Values{"mode": []string{"day"}, "filter": []string{"for_android"}}
		if g, e := r.Form, expectedForm; !reflect.DeepEqual(g, e) {
			t.Errorf("got form values %#v, want %#v", g, e)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture("fixtures/get_illust_ranking.json"))
	}))
	defer ts.Close()

	cli := NewClient(ClientConfig{
		TokenProvider: tp,
		BaseURL:       ts.URL,
	})

	ranking, err := cli.GetIllustRanking(NewGetIllustRankingParams().SetMode(RankingModeDay))
	if err != nil {
		t.Fatal(err)
	}

	if g, e := len(ranking.Illusts), 30; g != e {
		t.Errorf("got Illusts count %v, want %v", g, e)
	}

	expectedIllust00 := resp.GetIllustRankingIllust{
		ID:    64936066,
		Title: "♡",
		Type:  "illust",
		ImageURLs: map[string]string{
			"medium":        "https://i.pximg.net/c/540x540_70/img-master/img/2017/09/13/12/30/00/64936066_p0_master1200.jpg",
			"large":         "https://i.pximg.net/c/600x1200_90/img-master/img/2017/09/13/12/30/00/64936066_p0_master1200.jpg",
			"square_medium": "https://i.pximg.net/c/360x360_70/img-master/img/2017/09/13/12/30/00/64936066_p0_square1200.jpg",
		},
		Caption:  "9/12 Happy birthday!! (・８・)",
		Restrict: 0,
		User: resp.GetIllustRankingIllustUser{
			ID:      6996493,
			Name:    "Lpip",
			Account: "lpmya",
			ProfileImageURLs: map[string]string{
				"medium": "https://i.pximg.net/user-profile/img/2017/01/27/04/05/23/12061814_44196f064c0064fe89fdb6e719df20fe_170.png",
			},
			IsFollowed: false,
		},
		Tags: []resp.GetIllustRankingIllustTag{
			{Name: "ラブライブ!"},
			{Name: "南ことり"},
			{Name: "南ことり生誕祭2017"},
			{Name: "ミナリンスキー"},
			{Name: "メイド"},
			{Name: "12年に1度のことりちゃん年"},
			{Name: "(・8・)"},
			{Name: "ラブライブ!10000users入り"},
		},
		Tools:       []string{"CLIP STUDIO PAINT"},
		CreateDate:  "2017-09-13T12:30:00+09:00",
		PageCount:   1,
		Width:       650,
		Height:      936,
		SanityLevel: 2,
		Series:      resp.GetIllustRankingIllustSeries{ID: 0, Title: ""},
		MetaSinglePage: map[string]string{
			"original_image_url": "https://i.pximg.net/img-original/img/2017/09/13/12/30/00/64936066_p0.png",
		},
		MetaPages:      []resp.GetIllustRankingIllustMetaPage{},
		TotalView:      59452,
		TotalBookmarks: 13233,
		IsBookmarked:   false,
		Visible:        true,
		IsMuted:        false,
	}
	if g, e := ranking.Illusts[0], expectedIllust00; !reflect.DeepEqual(g, e) {
		t.Errorf("got Illusts[0] %#v, want %#v", g, e)
	}

	expectedIllust28 := resp.GetIllustRankingIllust{
		ID:    64914849,
		Title: "ことりちゃんHappy birthday (・8・)♡",
		Type:  "illust",
		ImageURLs: map[string]string{
			"square_medium": "https://i.pximg.net/c/360x360_70/img-master/img/2017/09/12/00/00/02/64914849_p0_square1200.jpg",
			"medium":        "https://i.pximg.net/c/540x540_70/img-master/img/2017/09/12/00/00/02/64914849_p0_master1200.jpg",
			"large":         "https://i.pximg.net/c/600x1200_90/img-master/img/2017/09/12/00/00/02/64914849_p0_master1200.jpg",
		},
		Caption:  "ことりちゃんおめでちゅん(・8・)♡<br />仕事で忙しくてあんまり時間が取れないので８時間くらいでサラっと描きました<br /><br />『ラブライブリンガー！ＵＲ 総集編』は現在各委託店にて好評発売中<br />メロンブックス → <a href=\"http://goo.gl/5pZ2Gx\" target=\"_blank\">http://goo.gl/5pZ2Gx</a>\u3000とらのあな → <a href=\"http://goo.gl/VHqabu\" target=\"_blank\">http://goo.gl/VHqabu</a>",
		Restrict: 0,
		User: resp.GetIllustRankingIllustUser{
			ID:      144203,
			Name:    "北原朋萌｡",
			Account: "kitaharakobo",
			ProfileImageURLs: map[string]string{
				"medium": "https://i.pximg.net/user-profile/img/2017/01/29/13/40/40/12071292_188a092ee853ccef0e1ebb84a51ec1e8_170.jpg",
			},
			IsFollowed: false,
		},
		Tags: []resp.GetIllustRankingIllustTag{
			{Name: "ラブライブ!"},
			{Name: "南ことり"},
			{Name: "ことり式雪だるま"},
			{Name: "南ことり生誕祭2017"},
			{Name: "スクフェス"},
			{Name: "(・8・)"},
			{Name: "12年に1度のことりちゃん年"},
			{Name: "ラブライブ!500users入り"},
		},
		Tools:       []string{},
		CreateDate:  "2017-09-12T00:00:02+09:00",
		PageCount:   2,
		Width:       789,
		Height:      1200,
		SanityLevel: 4,
		Series: resp.GetIllustRankingIllustSeries{
			ID:    0,
			Title: "",
		},
		MetaSinglePage: map[string]string{},
		MetaPages: []resp.GetIllustRankingIllustMetaPage{
			{
				ImageURLs: map[string]string{
					"original":      "https://i.pximg.net/img-original/img/2017/09/12/00/00/02/64914849_p0.jpg",
					"square_medium": "https://i.pximg.net/c/360x360_70/img-master/img/2017/09/12/00/00/02/64914849_p0_square1200.jpg",
					"medium":        "https://i.pximg.net/c/540x540_70/img-master/img/2017/09/12/00/00/02/64914849_p0_master1200.jpg",
					"large":         "https://i.pximg.net/c/600x1200_90/img-master/img/2017/09/12/00/00/02/64914849_p0_master1200.jpg",
				},
			},
			{
				ImageURLs: map[string]string{
					"square_medium": "https://i.pximg.net/c/360x360_70/img-master/img/2017/09/12/00/00/02/64914849_p1_square1200.jpg",
					"medium":        "https://i.pximg.net/c/540x540_70/img-master/img/2017/09/12/00/00/02/64914849_p1_master1200.jpg",
					"large":         "https://i.pximg.net/c/600x1200_90/img-master/img/2017/09/12/00/00/02/64914849_p1_master1200.jpg",
					"original":      "https://i.pximg.net/img-original/img/2017/09/12/00/00/02/64914849_p1.jpg",
				},
			},
		},
		TotalView:      13411,
		TotalBookmarks: 923,
		IsBookmarked:   false,
		Visible:        true,
		IsMuted:        false,
	}
	if g, e := ranking.Illusts[28], expectedIllust28; !reflect.DeepEqual(g, e) {
		t.Errorf("got Illusts[28] %#v, want %#v", g, e)
	}

	if g, e := ranking.NextURL, "https://app-api.pixiv.net/v1/illust/ranking?filter=for_android&mode=day&offset=30"; g != e {
		t.Errorf("got NextURL %q, want %q", g, e)
	}
}

func TestClient_GetIllustRanking_NotFound(t *testing.T) {
	tp := &mockTokenProvider{token: "ATN7bmWC7Kg1OneEqSPa9GxKm1l1uVHa8cQQKme7BGY"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(fixture("fixtures/api_error.json"))
	}))
	defer ts.Close()

	cli := NewClient(ClientConfig{TokenProvider: tp, BaseURL: ts.URL})

	_, err := cli.GetIllustRanking(NewGetIllustRankingParams().SetMode(RankingModeDay))
	if err == nil {
		t.Fatalf("GetIllustRanking() should return an error if 404 response is received")
	}

	respErr, ok := err.(ErrAPI)
	if !ok {
		t.Fatalf("GetIllustRanking() should return an ErrAPI if 404 response is received")
	}

	if g, e := respErr.StatusCode, http.StatusNotFound; g != e {
		t.Errorf("got StatusCode %v, want %v", g, e)
	}

	apiErrBody, err := respErr.Decode()
	if err != nil {
		t.Fatal(err)
	}

	expectedAPIErrBody := resp.APIErrorBody{
		Error: resp.APIError{
			UserMessage:        "指定されたエンドポイントは存在しません",
			Message:            "",
			Reason:             "",
			UserMessageDetails: map[string]interface{}{},
		},
	}
	if g, e := apiErrBody, expectedAPIErrBody; !reflect.DeepEqual(g, e) {
		t.Errorf("got APIErrorBody %#v, want %#v", g, e)
	}
}
