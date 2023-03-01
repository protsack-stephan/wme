package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/protsack-stephan/wme/pkg/api"
	"github.com/protsack-stephan/wme/pkg/auth"
	"github.com/protsack-stephan/wme/schema/v2"
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

	// Getting the current time in UTC.
	dte := time.Now().UTC()
	// Creating a new request for the API client with a filter.
	req := &api.Request{
		Filters: []*api.Filter{
			{
				Field: "in_language.identifier",
				Value: "uk",
			},
		},
	}

	// Calling the GetBatches function of the API client to get the batches for the current time with the specified filter.
	bts, err := clt.GetBatches(ctx, &dte, req)

	// If error occurs while getting the batches, application will panic.
	if err != nil {
		log.Panic(err)
	}

	// Creating a temporary file to store the downloaded batch.
	tmf, err := os.CreateTemp("", "tmp.tar.gz")

	if err != nil {
		log.Panic(err)
	}

	// Closing the temporary file at the end of the execution.
	defer tmf.Close()

	// Downloading the batch with the specified batch identifier to the temporary file.
	// Same can be done for Snapshots using the DownloadSnapshot function.
	// Also to directly read the Snapshot or a Batch ReadBatch or ReadSnapshot function can be used.
	// PartSize and Concurrency configuration can be found inside the client settings.
	if err := clt.DownloadBatch(ctx, &dte, bts[0].Identifier, tmf); err != nil {
		log.Panic(err)
	}

	// Resetting the file pointer to the beginning of the temporary file.
	_, _ = tmf.Seek(0, 0)

	// Defining a callback function to be called for each article in the batch.
	cbk := func(art *schema.Article) error {
		log.Println(art.Name)
		log.Println(len(art.Abstract))
		log.Print(art.Event.Type)
		return nil
	}

	// Reading all the articles from the temporary file using the callback function.
	if err := clt.ReadAll(ctx, tmf, cbk); err != nil {
		log.Panic(err)
	}
}
