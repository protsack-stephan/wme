package main

import (
	"context"
	"encoding/json"
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
		ids, _ := json.Marshal(evt.Data)
		log.Println(string(ids))
	}

	if err := fhs.PageDelete(ctx, time.Now(), cb); err != nil {
		log.Panic(err)
	}
}
