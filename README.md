# Wikimedia Enterprise API examples

This is a repository of SDKs for wikimedia enterprise APIs. At the moment, it includes SDKs for authentication, realtime and on-demand APIs. It also includes examples for working with these APIs. 
Note that all of the code inside this repository is not production ready. Although large parts of the code are covered with unit tests, this code was not tested in real world applications, so use it at your own risk.

### Code samples

1. Simple [authentication](example/auth/) example.

1. [Application](example/edits/) to calculate realtime edits per article.

1. Connecting to the [realtime stream](example/firehose/).

1. Using [on-demand APIs](example/ondemand/) (article and projects look up)

1. Connecting [Realtime V2 Beta](example/realtime/).


### SDK(s)

1. [Firehose client.](pkg/firehose/)

1. [Authentication client.](pkg/auth/)

1. [On-Demand client.](pkg/ondemand/)

1. [Realtime V2 beta.](pkg/realtime/)
