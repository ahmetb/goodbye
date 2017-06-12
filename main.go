package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

const (
	defaultConfigFile      = "/etc/goodbye/config.json"
	defaultPollingInterval = time.Minute * 5
)

type config struct {
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessSecret"`
}

func main() {
	log := logger.WithPrefix(
		logger.NewSyncLogger(logger.NewLogfmtLogger(os.Stdout)),
		"time", logger.DefaultTimestampUTC)

	log.Log("msg", "reading configuration")
	auth, err := readConfig(configPath())
	if err != nil {
		log.Log("msg", "failed to read configuration", "error", err)
		os.Exit(1)
	}

	api, err := mkClient(auth)
	if err != nil {
		log.Log("msg", "failed to initialize api client", "error", err)
		os.Exit(1)
	}

	log.Log("msg", "retrieving user profile")
	me, err := api.GetSelf(nil)
	if err != nil {
		log.Log("msg", "failed to fetch user's own profile", "error", err)
		os.Exit(1)
	}
	log.Log("msg", "authenticated", "self", me.ScreenName, "id", me.IdStr)

	if err := scan(log, api, me); err != nil {
		log.Log("error", err)
		os.Exit(1)
	}
}

func readConfig(path string) (config, error) {
	var c config
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return c, errors.Wrap(err, "failed to read config file")
	}
	return c, errors.Wrap(json.Unmarshal(b, &c), "failed to parse config file")
}

func configPath() string {
	v := os.Getenv("GOODBYE_CONFIG_PATH")
	if v != "" {
		return v
	}
	return defaultConfigFile
}

func mkClient(c config) (*anaconda.TwitterApi, error) {
	if c.ConsumerKey == "" {
		return nil, errors.New("twitter: consumerKey is not set")
	}
	if c.ConsumerSecret == "" {
		return nil, errors.New("twitter: consumerSecret is not set")
	}
	if c.AccessToken == "" {
		return nil, errors.New("twitter: accessToken is not set")
	}
	if c.AccessTokenSecret == "" {
		return nil, errors.New("twitter: accessSecret is not set")
	}

	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	api := anaconda.NewTwitterApi(c.AccessToken, c.AccessTokenSecret)
	api.SetLogger(anaconda.BasicLogger)
	return api, nil
}

func pollingInterval() (time.Duration, error) {
	v := os.Getenv("GOODBYE_POLLING_INTERVAL")
	if v != "" {
		d, err := time.ParseDuration(v)
		return d, errors.Wrap(err, "failed to parse custom interval")
	}
	return defaultPollingInterval, nil
}
