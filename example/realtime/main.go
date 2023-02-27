package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/protsack-stephan/wme/pkg/auth"
	"github.com/protsack-stephan/wme/pkg/realtime"
	"github.com/protsack-stephan/wme/schema/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()
	ath := auth.NewClient()

	lgn, err := ath.Login(ctx, &auth.LoginRequest{
		Username: os.Getenv("WME_USERNAME"),
		Password: os.Getenv("WME_PASSWORD"),
	})

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		err := ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
			RefreshToken: lgn.RefreshToken,
		})

		if err != nil {
			log.Println(err)
		}
	}()

	rlt := realtime.NewClient()
	rlt.SetAccessToken(lgn.AccessToken)

	hdl := func(art *schema.Article) error {
		dta, _ := json.Marshal(art)
		log.Println(string(dta))
		return nil
	}

	arq := &realtime.ArticlesRequest{
		Since:  time.Now().UTC().Add(-1 * time.Hour),
		Fields: []string{"name", "abstract", "event.type"},
		Filters: []schema.Filter{
			{
				Field: "is_part_of.identifier",
				Value: "enwiki",
			},
			{
				Field: "event.type",
				Value: "update",
			},
			{
				Field: "namespace.identifier",
				Value: 0,
			},
		},
	}

	if err := rlt.Articles(ctx, arq, hdl); err != nil {
		log.Panic(err)
	}
}
