# Wikimedia Enterprise Authentication API SDK

This is very basic Auth API SDK that will allow you to quickly get started with WME authentication.


1. Login example:

    ```go
    ath := auth.NewClient()

    lgn, err := ath.Login(ctx, &auth.LoginRequest{
      Username: os.Getenv("WME_USERNAME"),
      Password: os.Getenv("WME_PASSWORD"),
    })

    if err != nil {
      log.Panic(err)
    }

    log.Println(*lgn)
    ```

1. Revoke token example:

    ```go
    ath := auth.NewClient()

    err := ath.RevokeToken(ctx, &auth.RevokeTokenRequest{
      RefreshToken: os.Getenv("WME_REFRESH_TOKEN"),
    })

    if err != nil {
      log.Panic(err)
    }
    ```

1. Refresh token example:

    ```go
    ath := auth.NewClient()

    rft, err := ath.RefreshToken(ctx, &auth.RefreshTokenRequest{
      Username:     os.Getenv("WME_USERNAME"),
      RefreshToken: os.Getenv("WME_REFRESH_TOKEN"),
    })

    if err != nil {
      log.Panic(err)
    }

    log.Println(*rft)
    ```
