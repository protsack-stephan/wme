package main

import (
	"context"
	"log"
	"os"

	"github.com/protsack-stephan/wme/pkg/api"
	"github.com/protsack-stephan/wme/pkg/auth"
)

func main() {
	// Creating a new background context for the application.
	ctx := context.Background()

	// Creating a new instance of the auth client.
	ath := auth.NewClient()

	// Calling the Login function of the auth client to login with credentials from environment variables.
	lgn, err := ath.Login(ctx, &auth.LoginRequest{
		Username: os.Getenv("WME_USERNAME"),
		Password: os.Getenv("WME_PASSWORD"),
	})

	if err != nil {
		log.Panic(err)
	}

	// Deferring the RevokeToken function of the auth client to revoke the token after the execution of the code.
	defer func() {
		err := ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
			RefreshToken: lgn.RefreshToken,
		})

		if err != nil {
			log.Println(err)
		}
	}()

	// Creating a new instance of the API client.
	clt := api.NewClient()
	// Setting the access token for the client using the login access token.
	clt.SetAccessToken(lgn.AccessToken)

	// Create api request to filter down to English wikipedia and namespace 0.
	arq := &api.Request{
		Fields: []string{"name", "abstract", "is_part_of.*"},
		Filters: []*api.Filter{
			{
				Field: "namespace.identifier",
				Value: 0,
			},
			{
				Field: "is_part_of.identifier",
				Value: "enwiki",
			},
		},
		Limit: 1,
	}

	// Calling On-demand API to get the result.
	ats, err := clt.GetArticles(ctx, "Earth", arq)

	if err != nil {
		log.Panic(err)
	}

	// Looping through results to show the output.
	for _, art := range ats {
		log.Println(art.Name)
		log.Println(art.Abstract)
	}
}
