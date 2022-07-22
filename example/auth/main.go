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
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	})

	if err != nil {
		log.Panic(err)
	}

	_, err = ath.RefreshToken(ctx, &auth.RefreshTokenRequest{
		Username:     os.Getenv("USERNAME"),
		RefreshToken: lgn.RefreshToken,
	})

	if err != nil {
		log.Panic(err)
	}

	err = ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
		RefreshToken: lgn.RefreshToken,
	})

	if err != nil {
		log.Panic(err)
	}
}
