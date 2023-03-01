package main

import (
	"context"
	"log"
	"os"

	"github.com/protsack-stephan/wme/pkg/api"
	"github.com/protsack-stephan/wme/pkg/auth"
	"github.com/protsack-stephan/wme/schema/v2"
)

func main() {
	ctx := context.Background()
	ath := auth.NewClient()

	lgn, err := ath.Login(ctx, &auth.LoginRequest{
		Username: os.Getenv("WME_USERNAME"),
		Password: os.Getenv("WME_PASSWORD"),
	})

	if err != nil {
		log.Panic(err)
	}

	defer ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
		RefreshToken: lgn.RefreshToken,
	})

	clt := api.NewClient()
	clt.SetAccessToken(lgn.AccessToken)

	arq := &api.Request{
		Fields:  []string{"name", "abstract", "event.*"},
		Filters: []*api.Filter{},
	}
	cbk := func(art *schema.Article) error {
		log.Println(art.Name)
		log.Println(len(art.Abstract))
		log.Print(art.Event.Type)
		return nil
	}

	if err := clt.StreamArticles(ctx, arq, cbk); err != nil {
		log.Panic(err)
	}
}
