# Wikimedia Enterprise API examples

This repository contains examples of how WME realtime API can be used in real world scenarios. Note that all of the code inside this repository is not production ready. Although large parts of the code are covered with unit tests, this code was not tested in real world applications, so use it at your own risk.

The repository also has examples of working with on-demand API.


### Code samples

1. Simple [authentication](example/auth/) example.

1. [Application](example/edits/) to calculate realtime edits per article.

1. Connecting to the [realtime stream](example/firehose/).

1. Using [on-demand APIs](example/ondemand/) (article and projects look up)


### SDK(s)

1. [Firehose client.](pkg/firehose/)

1. [Authentication client.](pkg/auth/)

1. [On-Demand client.](pkg/ondemand/)
