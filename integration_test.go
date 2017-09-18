// +build integration

package pixiv

import (
	"context"
	"image"
	_ "image/jpeg"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/search2d/go-pixiv/resp"
)

func TestIntegration_GetIllustRanking(t *testing.T) {
	ctx := context.Background()

	tokenProvider := NewOauthTokenProvider(OauthTokenProviderConfig{
		Credential: Credential{
			Username:     os.Getenv("USERNAME"),
			Password:     os.Getenv("PASSWORD"),
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	})

	client := NewClient(ClientConfig{
		TokenProvider: tokenProvider,
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
	})

	date := time.Date(2017, 9, 1, 0, 0, 0, 0, time.UTC)

	illusts := []resp.GetIllustRankingIllust{}

	initial, err := client.GetIllustRanking(
		ctx,
		NewGetIllustRankingParams().SetMode(RankingModeDay).SetDate(date),
	)
	if err != nil {
		t.Fatal(err)
	}

	illusts = append(illusts, initial.Illusts...)

	next := initial.NextURL

	for {
		if len(next) == 0 {
			break
		}

		ranking, err := client.GetIllustRankingNext(ctx, next)
		if err != nil {
			t.Fatal(err)
		}

		illusts = append(illusts, ranking.Illusts...)

		next = ranking.NextURL
	}

	for offset, illust := range illusts {
		t.Logf("[%d] ID:%d Title:%q", offset, illust.ID, illust.Title)
	}
}

func TestIntegration_GetIllustDetail(t *testing.T) {
	ctx := context.Background()

	tokenProvider := NewOauthTokenProvider(OauthTokenProviderConfig{
		Credential: Credential{
			Username:     os.Getenv("USERNAME"),
			Password:     os.Getenv("PASSWORD"),
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	})

	client := NewClient(ClientConfig{
		TokenProvider: tokenProvider,
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
	})

	illust, err := client.GetIllustDetail(ctx, NewGetIllustDetailParams().SetIllustID(62397682))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(illust)
}

func TestIntegration_SetDownloadHeaders(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://i.pximg.net/img-original/img/2008/10/14/00/34/39/1859785_p0.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}

	SetDownloadHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if g, e := res.StatusCode, http.StatusOK; g != e {
		t.Errorf("got status code %d, want %d", g, e)
	}

	_, ext, err := image.Decode(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if g, e := ext, "jpeg"; g != e {
		t.Errorf("got image format %s, want %s", g, e)
	}
}
