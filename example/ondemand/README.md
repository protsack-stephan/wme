# Wikimedia Enterprise On-Demand API example

Simple piece of code that shows how to:

- do an article look up in a specific project.
- get all the available projects.

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

1. Note that if you are running this from the root use this command:
   ```bash
   go run example/ondemand/main.go
   ```
