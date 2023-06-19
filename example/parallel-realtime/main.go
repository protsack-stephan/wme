package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/protsack-stephan/wme/pkg/auth"
	"github.com/protsack-stephan/wme/pkg/realtime"
	"github.com/protsack-stephan/wme/schema/v2"
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

	rlt := realtime.NewClient()
	rlt.SetAccessToken(lgn.AccessToken)

	hdl := func(art *schema.Article) error {
		dta, _ := json.Marshal(art)
		log.Println(string(dta))
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		arq := &realtime.ArticlesRequest{
			Fields: []string{"name", "event.*"},
			Parts:  []int{2, 1}, // For partitions 0 through 49. This will connect to partitions 5 through 14
		}

		if err := rlt.Articles(ctx, arq, hdl); err != nil {
			log.Panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		arq := &realtime.ArticlesRequest{
			Fields: []string{"name", "event.*"},
			Parts:  []int{0, 1}, // Will connect to partitions 0 through 9
			// API will pick offsets for relevant partitions (0 through 9); for other partitions API will return messages from the earliest available
			// API will ignore irrelevant partitions in offset.
			Offsets: map[int]int64{19: 1, 22: 1, 15: 1, 21: 1, 33: 1, 38: 1, 44: 1, 3: 1, 5: 1, 7: 1, 34: 1, 13: 1, 24: 1, 30: 1, 20: 1, 25: 1, 39: 1, 42: 1, 145: 1, 4: 1, 14: 1, 18: 1, 146: 1, 47: 1, 29: 1, 40: 1, 41: 1, 43: 1, 149: 1, 9: 1, 27: 1, 28: 1, 10: 1, 31: 1, 36: 1, 37: 1, 0: 1, 1: 1, 2: 1},
			Filters: []realtime.Filter{
				{
					Field: "is_part_of.identifier",
					Value: "enwiki",
				},
				{
					Field: "event.type",
					Value: "update",
				},
				{
					Field: "namespace.identifier",
					Value: 0,
				},
			},
		}

		if err := rlt.Articles(ctx, arq, hdl); err != nil {
			log.Panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		arq := &realtime.ArticlesRequest{
			Fields: []string{"name", "event.*"},
			Parts:  []int{0, 9}, // Will connect to partitions 0,1,2,3,4,45,46,47,48,49
			// Time-Offsets for partitions 0,1,2,3 and 46; for other partitions consume from earliest
			// 146 will get ignored
			SincePerPartition: map[int]time.Time{
				0:   time.Now().UTC().Add(-40 * time.Hour),
				1:   time.Now().UTC().Add(-1 * time.Hour),
				2:   time.Now().UTC().Add(-1 * time.Hour),
				3:   time.Now().UTC().Add(-1 * time.Hour),
				46:  time.Now().UTC().Add(-10 * time.Hour),
				146: time.Now().UTC().Add(-10 * time.Hour),
			},
			Filters: []realtime.Filter{
				{
					Field: "is_part_of.identifier",
					Value: "eswiki",
				},
				{
					Field: "namespace.identifier",
					Value: 0,
				},
			},
		}

		if err := rlt.Articles(ctx, arq, hdl); err != nil {
			log.Panic(err)
		}
	}()

	wg.Wait()
}
