package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/protsack-stephan/wme/pkg/auth"
	"github.com/protsack-stephan/wme/pkg/firehose"
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

	fhs := firehose.NewClient()
	fhs.SetAccessToken(lgn.AccessToken)

	cb := func(evt *firehose.Event) {
		link, _ := url.QueryUnescape(evt.Data.URL)
		log.Printf("name: %s, identifier: %d, url: %s, dt: %s",
			evt.Data.Name,
			evt.Data.Identifier,
			link,
			evt.ID[0].Dt.Format(time.RFC3339))
	}

	if err := fhs.PageUpdate(ctx, time.Now(), cb); err != nil {
		log.Panic(err)
	}
}
