package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ahmetb/goodbye/v4/pkg/twitter"

	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

const (
	defaultPollingInterval = time.Minute * 10
)

var (
	flDaemon   bool
	flRunOnce  bool
	flHTTPAddr string
	flFilePath string
	flInterval time.Duration

	outStream       = os.Stderr
	pollingInterval time.Duration
	log             logger.Logger
	idsFile         file
)

func init() {
	log = mkLogger(outStream)
	flag.BoolVar(&flDaemon, "daemon", false, "run as a continuous process")
	flag.BoolVar(&flRunOnce, "run-once", false, "run once, save the follower list to file and exit")
	flag.StringVar(&flHTTPAddr, "http-addr", "", "listen on http portrun on request (e.g. :8080), specify ENV to use $PORT")
	flag.DurationVar(&flInterval, "check-interval", defaultPollingInterval, "(-daemon only) customize the check interval")

	flag.StringVar(&flFilePath, "followers-file", "", "GCS object URL to persist follower ids (e.g. gs://<bucket>/<obj>)")
	flag.Parse()

	if !flDaemon && !flRunOnce && flHTTPAddr == "" {
		log.Log("error", "one of -run-once, -daemon or -http-addr must be set")
		os.Exit(1)
	}

	if flHTTPAddr == "ENV" {
		p := os.Getenv("PORT")
		if p == "" {
			log.Log("error", "specified -http-addr=ENV but $PORT env var is not set")
			os.Exit(1)
		}
		flHTTPAddr = ":" + p
	}

	if (flDaemon && (flRunOnce || flHTTPAddr != "")) ||
		(flRunOnce && (flDaemon || flHTTPAddr != "")) ||
		(flHTTPAddr != "" && (flDaemon || flRunOnce)) {
		log.Log("error", "one of -run-once, -daemon or -http-addr can be used")
		os.Exit(1)
	}

	if (flRunOnce || flHTTPAddr != "") && flFilePath == "" {
		log.Log("error", "-followers-file should be specified to store data when -run-once or -http-addr specified")
		os.Exit(1)
	}
	if v := os.Getenv("GOODBYE_POLLING_INTERVAL"); v != "" {
		duration, err := time.ParseDuration(v)
		if err != nil {
			log.Log("error", errors.Wrapf(err, "failed to parse custom interval"))
			os.Exit(1)
		}
		log.Log("severity", "debug", "message", "parsed custom polling interval", "value", duration)
		pollingInterval = duration
	}

	if flFilePath != "" {
		f, err := openGCSObject(flFilePath)
		if err != nil {
			log.Log("error", err)
			os.Exit(1)
		}
		idsFile = f
	}
}

func main() {
	client, me, err := mkTwitterClient()
	if err != nil {
		log.Log("error", err)
		os.Exit(1)
	}
	log.Log("msg", "authenticated to twitter", "screen_name", me.ScreenName, "id", me.IDStr)

	if flDaemon {
		err = runLoop(log, client, me, pollingInterval)
	} else if flRunOnce {
		err = runOnce(log, client, me, idsFile)
	} else if flHTTPAddr != "" {
		log.Log("msg", "starting server", "addr", flHTTPAddr)
		http.HandleFunc("/goodbye", func(w http.ResponseWriter, r *http.Request) {
			if err := runOnce(log, client, me, idsFile); err != nil {
				log.Log("error", err)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				log.Log("msg", "success")
			}
		})
		err = http.ListenAndServe(flHTTPAddr, nil)
	} else {
		err = errors.New("unhandled mode")
	}
	if err != nil {
		log.Log("error", err)
		os.Exit(1)
	}
	log.Log("msg", "done")
}

func mkTwitterClient() (client twitter.Twitter, self twitter.TwitterProfile, _ error) {
	var consumerKey, consumerSecret, accessToken, accessTokenSecret string
	m := map[string]*string{
		"CONSUMER_KEY":        &consumerKey,
		"CONSUMER_SECRET":     &consumerSecret,
		"ACCESS_TOKEN":        &accessToken,
		"ACCESS_TOKEN_SECRET": &accessTokenSecret,
	}
	for k, t := range m {
		if v := os.Getenv(k); v == "" {
			return client, self, errors.Errorf("%s environment variable not set", k)
		} else {
			*t = v
		}
	}
	client = twitter.NewGoTwitter(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	me, err := client.Self()
	return client, me, errors.Wrap(err, "failed to fetch user's own profile")
}

func mkLogger(w io.Writer) logger.Logger {
	return logger.WithPrefix(logger.NewSyncLogger(logger.NewLogfmtLogger(outStream)), "timestamp", logger.DefaultTimestampUTC)
}
