# Goodbye

**Goodbye** is a Go application, when someone unfollows you on Twitter,
it will Direct Message you their Twitter handle.

At its core, it’s a program that downloads your twitter follower list in certain
intervals, and compares it with the previous list to find out who has unfollowed
you.

It has the following execution modes:

- **`-daemon`**: Process continously runs and fetches the followers every 5
  minutes (or specified `$GOODBYE_POLLING_INTERVAL`). Doesn't need extra storage
  to store followers list (stores in-memory).

- **`-http-addr`:** Process listens on specified port and checks follower list
  on each request to `GET /goodbye`. Needs extra storage (currently a GCS
  bucket) to store follower list. (Suitable for serverless environments)

- **`-run-once:`** Process runs once and exits. Needs extra storage (currently
  a GCS bucket) to store follower list. (Suitable for cronjobs.)

## Setup

First, you will need a `config.json` file containing Twitter API credentials.

Create a Twitter application on https://dev.twitter.com and access tokens for
the application. Save these and pass as environment variables to the
application:

```sh
export CONSUMER_KEY=[...] \
       CONSUMER_SECRET=[...] \
       ACCESS_TOKEN=[...] \
       ACCESS_TOKEN_SECRET=[...]

./goodbye -daemon
```

You can compile the application with Go compiler, or use the Dockerfile to build
a container image.

### Run with Docker

Clone this repo on a Linux server with docker installed and build the Docker
image in the repository root directory:

```sh
docker build -t goodbye .
```

Then run the container and give it the arguments it needs.

```sh
docker run -d --restart=always \
    -e CONSUMER_KEY=[...] \
    -e CONSUMER_SECRET=[...] \
    -e ACCESS_TOKEN=[...] \
    -e ACCESS_TOKEN_SECRET=[...] \
    --name=goodbye-agent \
       goodbye -daemon
```

Check if it is running fine: `docker logs -f goodbye-agent`.

## Set up follower ID storage (for `-run-once` and `-http-addr` modes)

The `-daemon` mode stores the follower IDs in memory, however other modes need
external storage (currently only Google Cloud Storage is supported) to store
the follower IDs.

1. Create a new Google Cloud Storage bucket to store follower IDs:

       BUCKET_NAME=pick-a-name
       gsutil mb gs://$BUCKET_NAME

1. Upload an empty file named `ids` to GCS bucket

       touch ids
       gsutil cp ./ids gs://$BUCKET_NAME/ids

Use this `gs://my_bucket/ids` value as the `-followers-file` argument.

### Demo

![Goodbye in action](http://i.imgur.com/FQr9Qjl.png)

## Author

Copyright (c) 2013-2019, [Ahmet Alp Balkan](https://twitter.com/ahmetb)

Made in Bellevue,WA and Seattle, WA with love.
