package pkg

import (
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
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

func LoadConfigFromEnv() error {
	var c Config
	if err := yaml.UnmarshalStrict([]byte(os.Getenv("GH_APP_CONFIG")), &c); err != nil {
		return errors.Wrap(err, "failed parsing configuration file")
	}
	Cfg = &c
	return nil
}
