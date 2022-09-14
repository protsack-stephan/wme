package main

import (
	"context"
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
		log.Printf("'%s' with event '%s', in project '%s'", art.Name, art.Event.Type, art.IsPartOf.Identifier)
		return nil
	}

	arq := &realtime.ArticlesRequest{
		Since:  time.Now().UTC().Add(-1 * time.Hour),
		Fields: []string{"name", "event.*", "is_part_of.*"},
		Filters: []realtime.Filter{
			{
				Field: "is_part_of.identifier",
				Value: "enwiki",
			},
			{
				Field: "event.type",
				Value: "update",
			},
		},
	}

	if err := rlt.Articles(ctx, arq, hdl); err != nil {
		log.Panic(err)
	}
}
