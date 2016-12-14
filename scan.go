package main

import (
	"fmt"
	"time"

	"net/url"

	"github.com/ChimeraCoder/anaconda"
	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

func scan(log *logger.Context, api *anaconda.TwitterApi, self anaconda.User) error {
	d, err := pollingInterval()
	if err != nil {
		return errors.Wrap(err, "cannot read polling interval")
	}

	log.Log("msg", "starting to scan periodically", "interval", d)

	var prev []int64
	for c := time.Tick(d); ; <-c { // ticks immediately at start

		log.Log("msg", "retrieving followers")
		cur, err := getFollowerIDs(log, api)
		if err != nil {
			return errors.Wrap(err, "failed to fetch followers")
		}
		log.Log("msg", "retrieved followers list", "n", len(cur), "prev", len(prev))

		if prev != nil {
			unfollowers := diff(prev, cur)
			log.Log("diff", len(unfollowers))
			for _, uid := range unfollowers {
				u, err := api.GetUsersShowById(uid, nil)
				if err != nil {
					log.Log("msg", "failed to fetch unfollower profile", "uid", uid, "error", err)
					continue
				}

				log.Log("msg", "unfollower", "name", u.ScreenName, "id", u.IdStr)
				if err := sendDM(api, self.ScreenName, u.ScreenName); err != nil {
					return errors.Wrap(err, "failed to send direct message")
				}
			}
		} else {
			log.Log("msg", "storing followers")
		}
		prev = cur
		log.Log("msg", "waiting until next run")
	}
}

func getFollowerIDs(log *logger.Context, api *anaconda.TwitterApi) ([]int64, error) {
	var out []int64

	ch := api.GetFollowersIdsAll(url.Values{
		"count": []string{"5000"}, // maximize responses in a page
	})
	for page := range ch {
		if err := page.Error; err != nil {
			return nil, err
		}
		out = append(out, page.Ids...)
		log.Log("msg", "got page response", "n", len(page.Ids), "total", len(out))
	}
	return out, nil
}

// diff finds elements present in prev but not in cur.
func diff(prev, cur []int64) []int64 {
	m := make(map[int64]bool, len(prev))
	for _, v := range prev {
		m[v] = true
	}
	for _, v := range cur {
		delete(m, v)
	}
	out := make([]int64, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// sendDM sends the self a direct message indicating the user unfollowed.
func sendDM(api *anaconda.TwitterApi, selfScreenName, unfollowerScreenName string) error {
	msg := fmt.Sprintf("@%s unfollowed you.", unfollowerScreenName)
	_, err := api.PostDMToScreenName(msg, selfScreenName)
	return err
}
