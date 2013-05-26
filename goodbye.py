#!/usr/bin/python
# coding=utf-8
import time
import json
import tweepy

AUTH_CONFIG_FILE = "auth.config"
GOODBYE_MESSAGE = "@%s sorry to see you unfollowing, goodbye!"
POLL_INTERVAL_SECS = 60*5


def main():
    config = {}
    try:
        with open(AUTH_CONFIG_FILE) as config_file:
            config = json.loads(config_file.read())
    except IOError:
        print('Twitter authentication credentials not found.')
        print('Step 1: Create a Twitter app https://dev.twitter.com/apps/new')
        print('Step 2: Go to OAuth Tool tab, copy consumer keys:')
        config['CONSUMER_KEY'] = raw_input('Paste Consumer key: ')
        config['CONSUMER_SECRET'] = raw_input('Paste Consumer secret: ')

        print('Step 3: Go to Details tab, press button "create access token"')
        print('Step 4: Copy generated access tokens:')
        config['ACCESS_TOKEN'] = raw_input('Paste Access Token: ')
        config['ACCESS_SECRET'] = raw_input('Paste Access Token secret: ')

        with open(AUTH_CONFIG_FILE, 'w') as config_file:
            config_file.write(json.dumps(config))

    auth = tweepy.OAuthHandler(config['CONSUMER_KEY'],
                               config['CONSUMER_SECRET'])
    auth.set_access_token(config['ACCESS_TOKEN'], config['ACCESS_SECRET'])
    api = tweepy.API(auth)

    me = api.me()
    print('Welcome, %s! You have %d followers.' % (me.screen_name,
                                                   me.followers_count))

    prev_follower_ids = get_followers_ids(api)
    print('This will check every %d seconds until someone unfollows...'
          % POLL_INTERVAL_SECS)

    while True:
        time.sleep(POLL_INTERVAL_SECS)

        new_follower_ids = get_followers_ids(api)
        diff = prev_follower_ids - new_follower_ids

        for unfollower_id in diff:
            user = api.get_user(unfollower_id)
            tweet = send_mention(api, user)
            print('@%s (%s), sent mention http://twitter.com/%s'
                  % (user.screen_name, user.name, user.screen_name))

        if diff:
            print('You have %d followers.' % len(new_follower_ids))

        prev_follower_ids = new_follower_ids


def get_followers_ids(api):
    follower_ids = set()
    for friend in tweepy.Cursor(api.followers_ids).items():
        follower_ids.add(friend)
    return follower_ids


def send_mention(api, user):
    return api.update_status(GOODBYE_MESSAGE % user.screen_name)

if __name__ == '__main__':
    main()
