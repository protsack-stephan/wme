package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/protsack-stephan/wme/pkg/api"
	"github.com/protsack-stephan/wme/pkg/auth"
	"github.com/protsack-stephan/wme/schema/v2"
)

func printrs(val interface{}, err error) {
	dta, _ := json.Marshal(val)

	if err != nil {
		log.Panic(err)
	}

	log.Println(string(dta))
}

func main() {
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

	clt := api.NewClient()
	clt.SetAccessToken(lgn.AccessToken)

	crq := &api.Request{
		Fields: []string{"name", "identifier"},
		Filters: []*schema.Filter{
			{
				Field: "identifier",
				Value: "wiki",
			},
		},
	}

	printrs(clt.GetCodes(ctx, crq))
	printrs(clt.GetCode(ctx, "wiki", crq))

	lrq := &api.Request{
		Fields: []string{"name", "identifier"},
		Filters: []*schema.Filter{
			{
				Field: "identifier",
				Value: "en",
			},
		},
	}

	printrs(clt.GetLanguages(ctx, lrq))
	printrs(clt.GetLanguage(ctx, "en", lrq))

	prq := &api.Request{
		Fields: []string{"name", "identifier"},
		Filters: []*schema.Filter{
			{
				Field: "identifier",
				Value: "enwiki",
			},
		},
	}

	printrs(clt.GetProjects(ctx, prq))
	printrs(clt.GetProject(ctx, "enwiki", prq))

	nrq := &api.Request{
		Fields: []string{"name", "identifier", "description"},
		Filters: []*schema.Filter{
			{
				Field: "identifier",
				Value: 0,
			},
		},
	}

	printrs(clt.GetNamespaces(ctx, nrq))
	printrs(clt.GetNamespace(ctx, 6, nrq))

	brq := &api.Request{
		Filters: []*schema.Filter{
			{
				Field: "in_language.identifier",
				Value: "en",
			},
			{
				Field: "namespace.identifier",
				Value: 0,
			},
		},
	}

	dte := time.Now().UTC()
	printrs(clt.GetBatches(ctx, &dte, brq))
	printrs(clt.GetBatch(ctx, &dte, "enwiki_namespace_0", brq))
	printrs(clt.HeadBatch(ctx, &dte, "enwiki_namespace_0"))

	cbk := func(art *schema.Article) error {
		log.Println(art.Name)
		log.Println(len(art.Abstract))
		return nil
	}

	if err := clt.ReadBatch(ctx, &dte, "afwiktionary_namespace_14", cbk); err != nil {
		log.Panic(err)
	}

	srq := &api.Request{
		Filters: []*schema.Filter{
			{
				Field: "in_language.identifier",
				Value: "en",
			},
			{
				Field: "namespace.identifier",
				Value: 0,
			},
		},
	}

	printrs(clt.GetSnapshots(ctx, srq))
	printrs(clt.GetSnapshot(ctx, "enwiki_namespace_0", srq))
	printrs(clt.HeadSnapshot(ctx, "enwiki_namespace_0"))

	if err := clt.ReadSnapshot(ctx, "afwikibooks_namespace_0", cbk); err != nil {
		log.Panic(err)
	}
}
