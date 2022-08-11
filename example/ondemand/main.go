package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/protsack-stephan/wme/pkg/auth"
	"github.com/protsack-stephan/wme/pkg/ondemand"
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

	odc := ondemand.NewClient()
	odc.SetAccessToken(lgn.AccessToken)

	// Article look up using SDK
	res, err := odc.Article(ctx, &ondemand.ArticleRequest{
		Project: "enwiki",
		Name:    "Steamship",
	})

	if err != nil {
		log.Println(err)
	}

	log.Printf("name: %s, identifier: %d\n wikitext: %s",
		res.Name,
		res.Identifier,
		res.ArticleBody.Wikitext,
	)

	// Projects look up using SDK
	prs, err := odc.Projects(ctx)

	if err != nil {
		log.Println(err)
	}

	fmt.Println("\n\nTotal number of projects: ", len(prs))
	fmt.Println("Names : identifier pairs of all the projects:")

	for _, prj := range prs {
		fmt.Println(prj.Name, " : ", prj.Identifier)
	}
}
