package common

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

//Config stores global configuration loaded from json file
type Configuration struct {
	Database       string `split_words:"true" default:"db-prod.yml"`
	DatabaseHost   string `split_words:"true" default:"175.106.13.14"`
	DatabaseUser   string `split_words:"true" default:"root"`
	DatabasePass   string `split_words:"true" default:"quantum"`
	DatabaseSchema string `split_words:"true" default:"quantumdns"`
	RedisHost      string `split_words:"true" default:"localhost"`
	RedisPort      string `split_words:"true" default:"6379"`
	Interval       uint64 `split_words:"true" default:"5"`
}

var Config Configuration

func LoadConfig() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
	})
	if err := envconfig.Process("CRON ACTIVITY", &Config); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	logrus.Info("Loaded configs: ", Config)
}
