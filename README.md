# Goodbye

**Goodbye** is a Go application, when someone unfollows you on Twitter,
it will Direct Message you their Twitter handle.

It can be run as a daemon that checks your account for unfollowers every 5
minutes (with Docker support), or as a serverless function on Google Cloud
Functions (GCF) with periodic invocation through Google Cloud Scheduler.

## Installation

First, you will need a `config.json` file containing Twitter API credentials.

Create a Twitter application on https://dev.twitter.com and access tokens for
the application. Then, create a `config.json` file with in the format:

```json
{
  "consumerKey": "<value>",
  "consumerSecret": "<value>",
  "accessToken": "<value>",
  "accessSecret": "<value>"
}
```

You have two options:

1. Run as a Docker container (requires an always-on Linux machine)
2. Run as a serverless function (requires Google Cloud Functions)

### Option 1: Installation with Docker

Clone this repo on a Linux server with docker installed and build the Docker
image in the repository root directory:

```sh
docker build -t goodbye .
```

Then run the container (specify path to `config.json` in `-v` argument):

```sh
docker run -d --restart=always \
    -v /your/path/to/config.json:/etc/goodbye/config.json \
    --name=goodbye-agent \
    goodbye
```

Check if it is running fine: `docker logs -f goodbye-agent`.

You can use `-e KEY=VALUE` format to "docker run" command to customize some
parameters through environment variables:

* `GOODBYE_CONFIG_PATH` path to the config file (defaults to
  `/etc/goodbye/config.json`)
* `GOODBYE_POLLING_INTERVAL` API polling interval duration in Go time.Duration
  format (defaults to `5m`)

### Option 2: Run as a serverless function on Google Cloud Functions

This is how I run it (recommended!) and it costs nearly 0$/month.

1. Clone this repo and navigate to `cmd/gcf` directory.
1. (Optional) If you have `go` installed, run `go build` here to see if it
   builds without any error messages.
1. Copy your `config.json` file here.

1. Create a new Google Cloud Storage bucket to store follower IDs:

       BUCKET_NAME=pick-a-name
       gsutil mb gs://$BUCKET_NAME

1. Upload an empty file named `ids` to GCS bucket

       touch ids
       gsutil cp ./ids gs://$BUCKET_NAME/ids

1. View your function's details on Google Cloud Console, note its Service
   Account field.

1. Use `gcloud` to give service account of the GCF app permissions on the GCS
   bucket:

       gsutil iam ch serviceAccount:"$(gcloud config get-value core/project)"@appspot.gserviceaccount.com:objectAdmin gs://$BUCKET_NAME

1. Use `gcloud` in this directory to create a function (change the
   `YOUR_BUCKET_NAME` occurrence below):

       gcloud alpha functions deploy goodbye \
         --memory 128MB \
         --trigger-http \
         --region us-central1 \
         --entry-point GoodbyeHandler \
         --runtime go111 \
         --set-env-vars GCS_BUCKET=$BUCKET_NAME,GCS_OBJECT=ids,GOODBYE_CONFIG_PATH=config.json

1. Visit the function's trigger URL and you should see an OK response.

1. Create a Google Cloud Scheduler job every 10 minutes to call the endpoint of
   your function.

### How it works

This program runs periodically download the list of your followers and compare
it with the previous version to see who is no longer there (meaning, unfollowed
you). Then it sends a DM to your with their twitter @username.

If you followed Option 1 (docker container), the machine should be running all
the time so the program can check the follower list every 5 minutes.

### Demo

![Goodbye in action](http://i.imgur.com/FQr9Qjl.png)

## Author

Copyright (c) 2013-2019, [Ahmet Alp Balkan](https://twitter.com/ahmetb)

Made in Bellevue,WA and Seattle, WA with love.
