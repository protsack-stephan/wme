# Wikimedia Enterprise New SDK examples

Examples for the SDK, that include [On-demand](ondemand/), [Realtime](streaming/) and [Realtime Batch](download/) examples.

### Getting started

1. First just export environment variables:

   ```bash
   export WME_USERNAME="superuser";
   export WME_PASSWORD="secret";
   ```

1. Then you should be able to just run:

   ```bash
   go run main.go
   ```

1. Note that if you are running this from the root use this command (in case you are calling `ondemand` example):

   ```bash
   go run example/api/ondemand/main.go
   ```
