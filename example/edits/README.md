# Wikimedia Enterprise Firehose Application example

Example of an app that will count number of edits per article  just to showcase the abilities of the API.

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

1. Note that if you are running this from the root use (keep it running for a while to fill the database with data):

    ```bash
    go run example/edits/main.go
    ```

1. After that `articles.db` file will be created. You'll need [SQLite CLI](https://www.sqlite.org/cli.html) or some other `SQLite` client installed to run the following queries on newly created file:
    
    * List of [top 10](list.sql) edited articles.
    * [Number of edits](language.sql) per language.
    * [Number of edits](project.sql) per project.
