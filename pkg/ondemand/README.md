# Wikimedia Enterprise On-Demand API SDK

The SDK is for the following APIs:
i) Article lookup - look up an article of a specific project.

ii) Available projects - get list of available projects with their identifier, language, url, etc.


### Getting started

1. First, create an ondemand client and associate an access token to it:

    ```go
  	ctx := context.Background()
	ath := auth.NewClient()

	lgn, err := ath.Login(ctx, &auth.LoginRequest{
		Username: os.Getenv("WME_USERNAME"),
		Password: os.Getenv("WME_PASSWORD"),
	})
	if err != nil {
		log.Panic(err)
	}
	
  	od := ondemand.NewClient()
	od.SetAccessToken(lgn.AccessToken)
    ```

2. Article look up example:

  ```go
  	req := &ondemand.ArticleRequest{
		Project: "enwiki",
		Name:    "Steamship",
	}

	res, err := od.Article(ctx, req)
	if err != nil {
		log.Println(err)
	}

	log.Printf("name: %s, identifier: %d\n wikitext: %s",
		res.Name,
		res.Identifier,
		res.ArticleBody.Wikitext,
	)
  ```

3. Getting total number of available projects example:

```go
  	res, err := od.Projects(ctx)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Total number of projects: ", len(res))
```