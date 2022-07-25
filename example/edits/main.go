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
		Username: os.Getenv("WME_USERNAME"),
		Password: os.Getenv("WME_PASSWORD"),
	})

	if err != nil {
		log.Panic(err)
	}

	defer ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
		RefreshToken: lgn.RefreshToken,
	})

	if err != nil {
		log.Panic(err)
	}

	fhs := firehose.NewClient()
	fhs.SetAccessToken(lgn.AccessToken)

	cmr := firehose.NewConnectionManger()
	cmr.Add(&firehose.Connection{
		Since:   time.Now(),
		Stream:  fhs.PageUpdate,
		Handler: func(evt *firehose.Event) {},
	})
	cmr.Add(&firehose.Connection{
		Since:   time.Now(),
		Stream:  fhs.PageDelete,
		Handler: func(evt *firehose.Event) {},
	})

	errs := make(chan error, 10000)
	go cmr.Connect(ctx, errs)

	for err := range errs {
		log.Println(err)
	}
}
