# Goodbye

**Goodbye** is a Go application, when someone unfollows you on Twitter,
it will Direct Message you their Twitter handle.

## Installation: the Docker way

Create a Twitter application on https://dev.twitter.com and access tokens
for the application before you start. Then, create a `config.json` file with
the credentials:

```json
{
  "consumerKey": "<value>",
  "consumerSecret": "<value>",
  "accessToken": "<value>",
  "accessSecret": "<value>"
}
```

Clone this repo and build the Docker image:

    $ docker build -t goodbye .

Then run the container:

```
docker run -d --restart=always \
    -v /path/to/config.json:/etc/goodbye/config.json \
    --name=goodbye-agent \
    goodbye
```

Check if it is running fine: `docker logs -f goodbye-agent`.

## Installation: the hard way

1. Clone this repository, go to source directory
2. Install dependencies `pip install -r requirements.txt`
3. Run program `./goodbye.py`

### Obtaining Twitter API Credentials

Go to https://dev.twitter.com/, create an application and access tokens 
for it, then save those to a `config.json` described above.

### How it works

This program has to be running all the time to check list of your followers
and take the diff. Therefore, make sure the program stays up all the time.

## Author

Copyright (c) 2013-2016, [Ahmet Alp Balkan](http://ahmetalpbalkan.com)

Made in Bellevue, WA with love.
