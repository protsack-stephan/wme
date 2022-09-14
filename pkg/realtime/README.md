# Wikimedia Enterprise Realtime API SDK

Allows to quickly connect and start using WME Realtime API.

### Getting started

Connect to the stream:

  ```go
  rlt := realtime.NewClient()
  rlt.SetAccessToken(os.Getenv("WME_ACCESS_TOKEN"))

  err := rlt.Articles(context.Background(), nil, func(art *schema.Article) error {
    log.Printf("name: %s\n", art.Name)
    return nil
  })

  if err != nil {
    log.Panic(err)
  }
  ```
