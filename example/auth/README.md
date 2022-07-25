# Wikimedia Enterprise Authentication API example

Simple piece of code that shows a really basic way to login, revoke token and refresh token with WME Authentication API.

### Getting started:

1. First just export environment variables:

    ```bash
    export WME_USERNAME="superuser";
    export WME_PASSWORD="secret";
    ```

1. Then you should be able to just run:

    ```bash
    go run main.go
    ```

1. Note that if you are running this from the root use this command:
    ```bash
    go run example/auth/main.go
    ```
