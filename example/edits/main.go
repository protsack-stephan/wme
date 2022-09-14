package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/protsack-stephan/wme/pkg/auth"
	"github.com/protsack-stephan/wme/pkg/firehose"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewArticle(evt *firehose.Event) *Article {
	link, _ := url.QueryUnescape(evt.Data.URL)

	return &Article{
		URL:        link,
		Name:       evt.Data.Name,
		IsPartOf:   evt.Data.IsPartOf.Identifier,
		InLanguage: evt.Data.InLanguage.Identifier,
		Edits:      0,
	}
}

type Article struct {
	URL        string `gorm:"primarykey"`
	Name       string
	IsPartOf   string
	InLanguage string
	Edits      int
}

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

	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(sqlite.Open("articles.db"), &gorm.Config{})

	if err != nil {
		log.Panic(err)
	}

	if err := db.AutoMigrate(&Article{}); err != nil {
		log.Panic(err)
	}

	fhs := firehose.NewClient()
	fhs.SetAccessToken(lgn.AccessToken)

	cmr := firehose.NewConnectionManger()
	cmr.Add(&firehose.Connection{
		Since:  time.Now(),
		Stream: fhs.PageUpdate,
		Handler: func(evt *firehose.Event) {
			art := NewArticle(evt)

			if err := db.FirstOrCreate(art).Error; err != nil {
				log.Println(err)
				return
			}

			art.Edits += 1

			if err := db.Save(art).Error; err != nil {
				log.Println(err)
			}
		},
	})
	cmr.Add(&firehose.Connection{
		Since:  time.Now(),
		Stream: fhs.PageDelete,
		Handler: func(evt *firehose.Event) {
			if err := db.Delete(NewArticle(evt)).Error; err != nil {
				log.Println(err)
			}
		},
	})

	errs := make(chan error, 10000)
	go cmr.Connect(ctx, errs)

	for err := range errs {
		log.Println(err)
	}
}
