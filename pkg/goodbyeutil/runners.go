package goodbyeutil

import (
	"fmt"
	"time"

	"goodbye/pkg/twitter"

	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

// RunOnce compares the user IDs given in prev to the current follower IDs and
// sends a DM to the self about unfollowed users and runs
func RunOnce(log logger.Logger, prev []int64, api twitter.Twitter, self twitter.TwitterProfile) ([]int64, error) {
	cur, err := api.GetFollowerIDs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch followers")
	}
	log.Log("msg", "retrieved followers list", "n", len(cur), "prev", len(prev))
	if prev == nil {
		return cur, nil
	}
	unfollowers := diff(prev, cur)
	log.Log("diff", len(unfollowers))
	for _, uid := range unfollowers {
		u, err := api.GetUserByID(uid)
		if err != nil {
			log.Log("msg", "failed to fetch unfollower profile", "uid", uid, "error", err)
			continue
		}

		log.Log("msg", "unfollower", "name", u.ScreenName, "id", u.IDStr)
		if err := sendDM(log, api, self.IDStr, u.ScreenName); err != nil {
			return nil, errors.Wrap(err, "failed to send direct message")
		}
	}
	return cur, nil
}

// RunLoop executes RunOnce periodically with given interval while preserving
// the follower counts in memory.
func RunLoop(log logger.Logger, api twitter.Twitter, self twitter.TwitterProfile, interval time.Duration) error {
	var prev []int64
	for c := time.Tick(interval); ; <-c { // ticks immediately at start
		cur, err := RunOnce(log, prev, api, self)
		if err != nil {
			return err
		}
		prev = cur
		log.Log("msg", "waiting until next run", "followers", len(cur))
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
func sendDM(log logger.Logger, api twitter.Twitter, selfID, unfollowerScreenName string) error {
	log.Log("msg", "sending dm", "self_id", selfID)
	msg := fmt.Sprintf("@%s unfollowed you.", unfollowerScreenName)
	err := api.SendDM(selfID, msg)
	return errors.Wrap(err, "failed to send DM")
}
