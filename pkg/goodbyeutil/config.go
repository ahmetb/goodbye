package goodbyeutil

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"goodbye/pkg/twitter"

	"github.com/pkg/errors"
)

const (
	defaultConfigFile = "/etc/goodbye/config.json"
)

type config struct {
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessSecret"`
}

// GetConfig returns a twitter client and authenticates to the current user
// profile from the configuration in the environment.
func GetConfig() (twitter.Twitter, twitter.TwitterProfile, error) {
	var me twitter.TwitterProfile
	auth, err := readConfig(configPath())
	if err != nil {
		return nil, me, errors.Wrap(err, "failed to read configuration")
	}
	api, err := mkClient(auth)
	if err != nil {
		return nil, me, errors.Wrap(err, "failed to initialize api client")
	}
	me, err = api.Self()
	if err != nil {
		return nil, me, errors.Wrap(err, "failed to fetch user's own profile")
	}
	return api, me, nil
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

func mkClient(c config) (twitter.Twitter, error) {
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
	return twitter.NewGoTwitter(c.ConsumerKey, c.ConsumerSecret, c.AccessToken, c.AccessTokenSecret), nil
}
