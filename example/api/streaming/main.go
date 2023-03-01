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
	// Create a new background context.
	ctx := context.Background()

	// Create a new authentication client.
	ath := auth.NewClient()

	// Attempt to log in using the environment's WME_USERNAME and WME_PASSWORD.
	lgn, err := ath.Login(ctx, &auth.LoginRequest{
		Username: os.Getenv("WME_USERNAME"),
		Password: os.Getenv("WME_PASSWORD"),
	})

	if err != nil {
		log.Panic(err)
	}

	// Make sure to revoke the token once we're done.
	defer func() {
		err := ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
			RefreshToken: lgn.RefreshToken,
		})

		if err != nil {
			log.Println(err)
		}
	}()

	// Create a new API client and set the access token to the one we just obtained.
	clt := api.NewClient()
	clt.SetAccessToken(lgn.AccessToken)

	// Create a new API request for articles, specifying the fields we want to retrieve and the filters to apply.
	arq := &api.Request{
		Fields:  []string{"name", "abstract", "event.*"},
		Filters: []*api.Filter{},
	}

	// Define a callback function to handle each article retrieved.
	cbk := func(art *schema.Article) error {
		log.Println(art.Name)
		log.Println(len(art.Abstract))
		log.Print(art.Event.Type)
		return nil
	}

	// Stream the articles using the API client and the request we just defined.
	if err := clt.StreamArticles(ctx, arq, cbk); err != nil {
		log.Panic(err)
	}
}
