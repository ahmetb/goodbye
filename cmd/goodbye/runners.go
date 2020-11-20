package main

import (
	"fmt"
	"time"

	"github.com/ahmetb/goodbye/v4/pkg/twitter"

	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

// runLoop executes RunOnce periodically with given interval while preserving
// the follower counts in memory.
func runLoop(log logger.Logger, api twitter.Twitter, self twitter.TwitterProfile, interval time.Duration) error {
	var prev []int64
	log.Log("message", "starting loop", "interval", interval)
	for c := time.Tick(interval); ; <-c { // ticks immediately at start
		cur, err := run(log, prev, api, self)
		if err != nil {
			return err
		}
		prev = cur
		log.Log("message", "waiting until next run", "followers", len(cur))
	}
}

// runOnce executes the check loop once by reading the known follower list from
// the specified file and saves the reslts back to the specified file
func runOnce(log logger.Logger, api twitter.Twitter, self twitter.TwitterProfile, f file) error {
	log.Log("severity", "debug", "message", "loading previous IDs from file")
	prevIDs, err := loadIDs(f)
	if err != nil {
		return errors.Wrap(err, "failed to load current IDs")
	}

	curIDs, err := run(log, prevIDs, api, self)
	if err != nil {
		return errors.Wrap(err, "failed to run check")
	}

	err = saveIDs(f, curIDs)
	return errors.Wrap(err, "failed to save fetched list of followers")
}

func run(log logger.Logger, prev []int64, api twitter.Twitter, self twitter.TwitterProfile) ([]int64, error) {
	log.Log("severity", "debug", "message", "fetching follower list")
	cur, err := api.GetFollowerIDs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch followers")
	}
	log.Log("message", "retrieved followers list", "n", len(cur), "prev", len(prev))
	if prev == nil {
		return cur, nil
	}
	unfollowers := diff(prev, cur)
	log.Log("diff", len(unfollowers))
	for _, uid := range unfollowers {
		u, err := api.GetUserByID(uid)
		if err != nil {
			log.Log("message", "failed to fetch unfollower profile", "uid", uid, "error", err)
			continue
		}

		log.Log("message", "unfollower", "name", u.ScreenName, "id", u.IDStr)
		if err := sendDM(log, api, self.IDStr, u); err != nil {
			return nil, errors.Wrap(err, "failed to send direct message")
		}
	}
	return cur, nil
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
func sendDM(log logger.Logger, api twitter.Twitter,
	selfID string, unfollower twitter.TwitterProfile) error {
	log.Log("message", "sending dm", "self_id", selfID)
	msg := fmt.Sprintf("@%s unfollowed you. (%d/%d)",
		unfollower.ScreenName, unfollower.FollowingCount, unfollower.FollowerCount)
	err := api.SendDM(selfID, msg)
	return errors.Wrap(err, "failed to send DM")
}
