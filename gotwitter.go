package main

import (
	gotw "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	logger "github.com/go-kit/kit/log"
)

// 	goTwitter uses https://github.com/dghubble/go-twitter.
type goTwitter struct {
	log    logger.Logger
	client *gotw.Client
}

func newGoTwitter(log logger.Logger, c config) twitter {
	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := gotw.NewClient(httpClient)
	return &goTwitter{
		log:    log,
		client: client}
}

func (g *goTwitter) Self() (twitterProfile, error) {
	u, _, err := g.client.Accounts.VerifyCredentials(&gotw.AccountVerifyParams{})
	if err != nil {
		return twitterProfile{}, err
	}
	return twitterProfile{
		screenName: u.ScreenName,
		idStr:      u.IDStr}, nil
}

func (g *goTwitter) GetFollowerIDs() ([]int64, error) {
	var out []int64
	var cursor int64

	page := 0
	for {
		g.log.Log("msg", "fetching follower ids", "page", page, "cursor", cursor)
		ids, _, err := g.client.Followers.IDs(&gotw.FollowerIDParams{
			Count:  5000,
			Cursor: cursor})
		if err != nil {
			return nil, err
		}

		out = append(out, ids.IDs...)
		cursor = ids.NextCursor
		g.log.Log("msg", "fetched follower ids", "page", page, "count", len(ids.IDs), "next_cursor", cursor)
		if cursor == 0 {
			break
		}
		page++
	}
	g.log.Log("msg", "fetched all follower ids", "count", len(out))
	return out, nil
}

func (g *goTwitter) SendDM(to, message string) error {
	dm, _, err := g.client.DirectMessages.EventsNew(&gotw.DirectMessageEventsNewParams{
		Event: &gotw.DirectMessageEvent{
			Type: "message_create",
			Message: &gotw.DirectMessageEventMessage{
				Target: &gotw.DirectMessageTarget{
					RecipientID: to,
				},
				Data: &gotw.DirectMessageData{
					Text: message,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	g.log.Log("msg", "sent DM", "id", dm.ID)
	return nil
}

func (g *goTwitter) GetUserByID(id int64) (twitterProfile, error) {
	u, _, err := g.client.Users.Show(&gotw.UserShowParams{UserID: id})
	if err != nil {
		return twitterProfile{}, err
	}
	return twitterProfile{
		screenName: u.ScreenName,
		idStr:      u.IDStr}, nil
}
