#!/usr/bin/python
# coding=utf-8
import time
import json
import random
import tweepy
import os

AUTH_CONFIG_FILE = "auth.config"
POLL_INTERVAL_SECS = 60*5 

def main():
    api=config_from_env()
    if not api:
        print('Environment variables are not set, falling back to configuration file {0}'.format(AUTH_CONFIG_FILE))
        api=config_from_file()
    else:
        print('Picked up configuration from environment variables.')

    me = api.me()
    print('Welcome, %s! You have %d followers.' % (me.screen_name,
                                                   me.followers_count))

    prev_follower_ids = get_followers_ids(api)
    print('This will check every %d seconds until someone unfollows...'
          % POLL_INTERVAL_SECS)

    mentioned = set()

    while True:
        time.sleep(POLL_INTERVAL_SECS)

        new_follower_ids = get_followers_ids(api)

        diff = set()
        if new_follower_ids:
            diff = prev_follower_ids - new_follower_ids - mentioned
            prev_follower_ids = new_follower_ids

        for unfollower_id in diff:
            try:
                user = get_twitter_user(api, unfollower_id)

                if user:
                    tweet = send_dm(api, user, me.screen_name)
                    print('unfollow: @%s (%s) http://twitter.com/%s'
                          % (user.screen_name, user.name, user.screen_name))
                    mentioned.add(unfollower_id)
            except Exception as e:
                print 'Cannot send tweet.', e

        if diff:
            print('You have %d followers now.' % len(new_follower_ids))

def config_from_env():
    keys = ['CONSUMER_KEY', 'CONSUMER_SECRET', 'ACCESS_TOKEN', 'ACCESS_SECRET']
    for k in keys:
        if k not in os.environ:
	    print("${0} not found in environment.".format(k))
	    return
    return api_from_config(os.environ)

def config_from_file():
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
    return api_from_config(config)

def api_from_config(config):
    auth = tweepy.OAuthHandler(config['CONSUMER_KEY'],
                               config['CONSUMER_SECRET'])
    auth.set_access_token(config['ACCESS_TOKEN'], config['ACCESS_SECRET'])
    return tweepy.API(auth)

def get_twitter_user(api, user_id):
    try:
        return api.get_user(user_id)
    except Exception as e:
        print 'Cannot fetch user with id {0}: {1}'.format(user_id, e)


def get_followers_ids(api):
    try:
        return set(tweepy.Cursor(api.followers_ids).items())
    except Exception as e:
        print 'Cannot fetch followers:', e


def get_random(arr):
    if arr:
        random.shuffle(arr)
        return arr[0]


def send_dm(api, user, me):
    tweet = 'D {0} @{1} unfollowed you.'.format(me, user.screen_name)
    return api.update_status(tweet)

if __name__ == '__main__':
    main()
