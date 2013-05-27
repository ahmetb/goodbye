#!/usr/bin/python
# coding=utf-8
import time
import json
import random
import tweepy

AUTH_CONFIG_FILE = "auth.config"
GOODBYE_MESSAGES_FILE = "messages.txt"
POLL_INTERVAL_SECS = 5 * 60

def main():
    config = {}
    messages = []

    try:
        with open(GOODBYE_MESSAGES_FILE) as messages_file:
            messages = [m.strip() for m in messages_file.readlines()]
    except IOError: pass
    if not messages:
        print('No messages found in %s file' % GOODBYE_MESSAGES_FILE)
        return

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
            tweet = send_mention(api, user, messages)
            print('@%s (%s), sent mention http://twitter.com/%s'
                  % (user.screen_name, user.name, user.screen_name))

        if diff:
            print('You have %d followers.' % len(new_follower_ids))

        prev_follower_ids = new_follower_ids


def get_followers_ids(api):
    return set(tweepy.Cursor(api.followers_ids).items())


def get_random(arr):
    if arr:
        random.shuffle(arr)
        return arr[0]

def send_mention(api, user, messages):
    content = get_random(messages)
    tweet = '@{0} {1}'.format(user.screen_name, content)
    return api.update_status(tweet)

if __name__ == '__main__':
    main()
