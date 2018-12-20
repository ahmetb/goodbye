package main

import (
	"os"
	"time"

	"github.com/ahmetb/goodbye/v3/pkg/goodbyeutil"

	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

const (
	defaultPollingInterval = time.Minute * 5
)

func run(log logger.Logger) error {
	var interval time.Duration
	if v := os.Getenv("GOODBYE_POLLING_INTERVAL"); v != "" {
		duration, err := time.ParseDuration(v)
		if err != nil {
			errors.Wrapf(err, "failed to parse custom interval")
		}
		interval = duration
	}
	api, me, err := goodbyeutil.GetConfig()
	if err != nil {
		return errors.Wrap(err, "failed to initialize")
	}
	log.Log("message", "authenticated", "screen_name", me.ScreenName, "id", me.IDStr)
	log.Log("message", "starting to run periodically", "interval", interval)
	err = goodbyeutil.RunLoop(log, api, me, interval)
	return errors.Wrap(err, "run loop terminated with error")
}

func main() {
	log := logger.WithPrefix(
		logger.NewSyncLogger(logger.NewLogfmtLogger(os.Stdout)), "timestamp", logger.DefaultTimestampUTC)

	if err := run(log); err != nil {
		log.Log("severity", "error", "error", err)
		os.Exit(1)
	}
}
