package pkg

import (
	"os"
	"strconv"

	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server             HTTPConfig       `yaml:"server"`
	Github             githubapp.Config `yaml:"github"`
	ExpectedPusherName string           `yaml:"expected_pusher_name"`
}

type HTTPConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

var Cfg *Config

func LoadConfig(path string) error {
	var c Config

	bytes, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "failed reading server config file: %s", path)
	}

	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		return errors.Wrap(err, "failed parsing configuration file")
	}
	Cfg = &c
	return nil
}

func LoadConfigFromEvn() error {
	var c Config
	c.Github.App.PrivateKey = os.Getenv("GH_APP_PRIVATE_KEY_PEM")
	integrationId, err := strconv.Atoi(os.Getenv("GH_APP_INTEGRATION_ID"))
	if err != nil {
		return err
	}
	c.Github.App.IntegrationID = int64(integrationId)
	c.Github.App.WebhookSecret = os.Getenv("GH_APP_WEBHOOK_SECRET")
	Cfg = &c
	return nil
}
