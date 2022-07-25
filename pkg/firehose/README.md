# Wikimedia Enterprise Firehose (realtime) API SDK

Allows to quickly connect and start using WME firehose (realtime) API.

### Getting started

1. Connect to the stream:

    ```go
    fhs := firehose.NewClient()
    fhs.SetAccessToken(os.Getenv("WME_ACCESS_TOKEN"))

    cb := func(evt *firehose.Event) {
      link, _ := url.QueryUnescape(evt.Data.URL)
      log.Printf("name: %s, identifier: %d, url: %s, dt: %s",
        evt.Data.Name,
        evt.Data.Identifier,
        link,
        evt.ID[0].Dt.Format(time.RFC3339))
    }

    if err := fhs.PageUpdate(context.Background(), time.Now(), cb); err != nil {
      log.Panic(err)
    }
    ```
