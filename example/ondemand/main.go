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
	// fmt.Println("login response")
	// spew.Dump(lgn)

	// defer ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
	// 	RefreshToken: lgn.RefreshToken,
	// })
	od := ondemand.NewClient()
	od.SetAccessToken(lgn.AccessToken)
	// od.SetAccessToken("eyJraWQiOiJzeVNnS1JaZWdwcDFlSGZEYnlsR2YrTnBjVmVXUDZJNGJlSFpOWjBDZVdrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1ZDI3NjQyZi02YTE2LTQyYjEtOGE3ZC1hYWQ2NmZiNGUyMTIiLCJjb2duaXRvOmdyb3VwcyI6WyJncm91cF8xIl0sImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy1lYXN0LTEuYW1hem9uYXdzLmNvbVwvdXMtZWFzdC0xX0tiNW5ZZDN6dSIsImNsaWVudF9pZCI6IjY0MXU0aTdncHR1ZmZzc2w0bTlvYXR2NHU5Iiwib3JpZ2luX2p0aSI6ImY1NGQ0YzEzLWY4NzEtNGY2MC1iZDVlLTk4MzhkYzI3NmJjZiIsImV2ZW50X2lkIjoiNzBkYTQxODEtNzhmZC00MmU3LTlmNWItMzMzYTMzNWE2OTA3IiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTY1OTQ3MTg1MiwiZXhwIjoxNjU5NTU4MjUyLCJpYXQiOjE2NTk0NzE4NTIsImp0aSI6IjEyYWE1OTM5LTJjYjEtNGU4ZS1hMmE2LWQ0ZGRlMWM2OTBhZCIsInVzZXJuYW1lIjoicmF0YXRvc2tyIn0.CzABtsoBusEkrMq4TksjGnRzagn8eb46uV_o_lGt3bltc_y92JKRTo8gD5-yfXa2n0VF2ciI7-4_2kudRFaLoHOjBIKaZRmNXjsK3VoImlcsbRCHtMBymt93KUXbBCZlqP6i7jBx4bLQ4i6TCHo9RBppD4mjn4--Qn8azwnmQ_Oj2aphj_pHfc_XsEUACWIzLXyEzo2oHFMTw4kAnvAqkx9fBqKy3O6PniFRVvr5VtW08dAwO22zWsyw_qZlwhQ_-Lmk0QiJZhT6RsqoqBqttN2pOKDUe7OxkAJyX-51VAX1x86xCSwzbQsx35oec3x4umswgX4GgDLO9p72H0njMQ")
	req := &ondemand.ArticleRequest{
		Project: "enwiki",
		Name:    "Steamship",
	}

	res, err := od.Article(ctx, req)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(" response")
	log.Printf("name: %s, identifier: %d\n wikitext: %s",
		res.Name,
		res.Identifier,
		res.ArticleBody.Wikitext,
	)

	prjs, err := od.Projects(ctx)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Total number of projects: ", len(prjs))
	fmt.Println("Names : identifier of all the projects are as follows:")
	for _, project := range prjs {
		fmt.Println(project.Name, " : ", project.Identifier)
	}
}
