package twitter

import (
	gotw "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// GoTwitter uses https://github.com/dghubble/go-twitter.
type GoTwitter struct {
	client *gotw.Client
}

// NewGoTwitter initializes a go-twitter client.
func NewGoTwitter(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *GoTwitter {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := gotw.NewClient(httpClient)
	return &GoTwitter{client: client}
}

func (g *GoTwitter) Self() (TwitterProfile, error) {
	u, _, err := g.client.Accounts.VerifyCredentials(&gotw.AccountVerifyParams{})
	if err != nil {
		return TwitterProfile{}, err
	}
	return TwitterProfile{
		ScreenName: u.ScreenName,
		IDStr:      u.IDStr}, nil
}

func (g *GoTwitter) GetFollowerIDs() ([]int64, error) {
	var out []int64
	var cursor int64

	page := 0
	for {
		// g.log.Log("msg", "fetching follower ids", "page", page, "cursor", cursor)
		ids, _, err := g.client.Followers.IDs(&gotw.FollowerIDParams{
			Count:  5000,
			Cursor: cursor})
		if err != nil {
			return nil, err
		}

		out = append(out, ids.IDs...)
		cursor = ids.NextCursor
		// g.log.Log("msg", "fetched follower ids", "page", page, "count", len(ids.IDs), "next_cursor", cursor)
		if cursor == 0 {
			break
		}
		page++
	}
	// g.log.Log("msg", "fetched all follower ids", "count", len(out))
	return out, nil
}

func (g *GoTwitter) SendDM(to, message string) error {
	_, _, err := g.client.DirectMessages.EventsNew(&gotw.DirectMessageEventsNewParams{
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
	// g.log.Log("msg", "sent DM", "id", dm.ID)
	return nil
}

func (g *GoTwitter) GetUserByID(id int64) (TwitterProfile, error) {
	u, _, err := g.client.Users.Show(&gotw.UserShowParams{UserID: id})
	if err != nil {
		return TwitterProfile{}, err
	}
	return TwitterProfile{
		ScreenName: u.ScreenName,
		IDStr:      u.IDStr}, nil
}
