// +build integration

package pixiv

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/search2d/go-pixiv/resp"
)

func TestIntegration_GetIllustRanking(t *testing.T) {
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

		ranking, err := client.GetIllustRankingNext(next)
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
