# Wikimedia Enterprise APIs SDK

Intended to jump start you with the using of the API(s).

## Getting started

To create the API client you can use the following code snippet:

```go
clt := api.NewClient()
```

If you need to pass a custom options to the client you can use:

```go

clt := api.NewClient(func(clt *api.Client) {
  // redefine the client properties here
  clt.UserAgent = "my agent"
})
```

Also, to use the api you need to set the access token:

```go
clt := api.NewClient()
clt.SetAccessToken("my_token")
```

Please refer to the [interface](api.go#L59) definitions to see the full list of APIs.
