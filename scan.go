package main

import (
	"fmt"
	"time"

	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

func scan(log logger.Logger, api twitter, self twitterProfile) error {
	d, err := pollingInterval()
	if err != nil {
		return errors.Wrap(err, "cannot read polling interval")
	}

	log.Log("msg", "starting to scan periodically", "interval", d)

	var prev []int64
	for c := time.Tick(d); ; <-c { // ticks immediately at start

		log.Log("msg", "retrieving followers")
		cur, err := api.GetFollowerIDs()
		if err != nil {
			return errors.Wrap(err, "failed to fetch followers")
		}
		log.Log("msg", "retrieved followers list", "n", len(cur), "prev", len(prev))

		if prev != nil {
			unfollowers := diff(prev, cur)
			log.Log("diff", len(unfollowers))
			for _, uid := range unfollowers {
				u, err := api.GetUserByID(uid)
				if err != nil {
					log.Log("msg", "failed to fetch unfollower profile", "uid", uid, "error", err)
					continue
				}

				log.Log("msg", "unfollower", "name", u.screenName, "id", u.idStr)
				if err := sendDM(log, api, self.idStr, u.screenName); err != nil {
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
func sendDM(log logger.Logger, api twitter, selfID, unfollowerScreenName string) error {
	log.Log("msg", "sending dm", "self_id", selfID)
	msg := fmt.Sprintf("@%s unfollowed you.", unfollowerScreenName)
	err := api.SendDM(selfID, msg)
	return errors.Wrap(err, "failed to send DM")
}
