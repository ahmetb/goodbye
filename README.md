# Goodbye

> It is a Turkish tradition to throw water behind a vehicle when the loved ones are setting out on a long journey.
> The tradition itself serves to express a wish —that the long journey will go smoothly, without mishap. As smooth as water.

![turkish tradition: throwing water](https://github.com/ahmetalpbalkan/goodbye/raw/master/img/promo.gif)

**Goodbye** is a Python application, when someone unfollows you on Twitter, it will say them goodbye as long at it is running.

### Installation

1. Clone this repository, go to source directory
2. Install dependencies `pip install -r requirements.txt`
3. Run program `./goodbye.py`

When the program is run for the first time, it will ask for configuration regarding Twitter app  you will create and authentication tokens only once.

Go to https://dev.twitter.com/ , create an application and access tokens for it and paste them to the program. Example:

![](http://i.imgur.com/CQPgJaM.png)

### How it works

This program has to be running all the time to say goodbye to your unfollowers.

Start program with `./goodbye.py` and keep it running. You can host this on your server, Raspberry Pi or some always running computer.

Just terminate program with `Ctrl+C` to stop tracking.

### Advanced Configuration

In `goodbye.py` file, there are configuration keys:

    POLL_INTERVAL_SECS = 60*5

You can adjust Twitter API polling interval as well, e.g. defaults to 5 minutes above.

You can customize random goodbye messages by editing `messages.txt` by writing
one message per line. Try to keep each message less than 110 characters.

## Demo

Try following me [`@ahmetalpbalkan`](http://twitter.com/ahmetalpbalkan) on Twitter and unfollow 10 minutes later (works only if I keep this program running).

## Author

Copyright (c) 2013, [Ahmet Alp Balkan](http://ahmetalpbalkan.com)

Made in Bellevue, WA with love.
