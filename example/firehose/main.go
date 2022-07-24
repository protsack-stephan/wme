package main

import (
	"context"
	"log"
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
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	})

	if err != nil {
		log.Panic(err)
	}

	defer ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
		RefreshToken: lgn.RefreshToken,
	})

	fhs := firehose.NewClient()
	fhs.SetAccessToken(lgn.AccessToken)

	cb := func(evt *firehose.Event) {
		log.Println(*evt)
		log.Println(*evt.Data)
	}

	if err := fhs.PageUpdate(ctx, time.Now(), cb); err != nil {
		log.Panic(err)
	}
}
