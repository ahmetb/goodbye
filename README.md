# Goodbye

**Goodbye** is a Python application, when someone unfollows you on Twitter,
it will Direct Message you their twitter handle.

### Installation the Docker way

Make sure you have a Twitter OAuth(xAuth) credentials before you start.
If not, run the program manually once and copy them from `auth.config` file
to the command below.

Clone this repo and build an image:

    $ docker build -t goodbye .

Then run the container:

```
docker run -d --restart=always \
    -e CONSUMER_KEY=<value> \
    -e CONSUMER_SECRET=<value> \
    -e ACCESS_TOKEN=<value> \
    -e ACCESS_SECRET=<value> \
    --name=goodbye-agent \
    goodbye
```

Check if it is running fine: `docker logs -f goodbye-agent`.

## Installation the hard way

1. Clone this repository, go to source directory
2. Install dependencies `pip install -r requirements.txt`
3. Run program `./goodbye.py`

When the program is run for the first time, it will ask for configuration regarding Twitter app  you will create and authentication tokens only once.

Go to https://dev.twitter.com/ , create an application and access tokens for it and paste them to the program. Example:

![](http://i.imgur.com/CQPgJaM.png)

### How it works

This program has to be running all the time to check list of your followers
and take the diff. Therefore, make sure the program stays up all the time.

### Advanced Configuration

In `goodbye.py` file, there are configuration keys:

    POLL_INTERVAL_SECS = 60*5

You can adjust Twitter API polling interval as well, e.g. defaults to 5 minutes above.

## Author

Copyright (c) 2013, [Ahmet Alp Balkan](http://ahmetalpbalkan.com)

Made in Bellevue, WA with love.
