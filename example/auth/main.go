package main

import (
	"context"
	"log"
	"os"

	"github.com/protsack-stephan/wme/pkg/auth"
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

	rft, err := ath.RefreshToken(ctx, &auth.RefreshTokenRequest{
		Username:     os.Getenv("WME_USERNAME"),
		RefreshToken: lgn.RefreshToken,
	})

	if err != nil {
		log.Panic(err)
	}

	log.Println(rft)

	err = ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
		RefreshToken: lgn.RefreshToken,
	})

	if err != nil {
		log.Panic(err)
	}
}
