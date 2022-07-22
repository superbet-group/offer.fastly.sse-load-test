package config

import "github.com/superbet-group/trading.commons/pkg/configlib"

// Config contains all app configurations
type Config struct {
	Version             string `default:"v.0.0.0"`
	AppEnvironment      string `split_words:"true" default:"local"`
	App                 configlib.AppConfig
	OfferHost           string `split_words:"true"`
	Monitoring          Monitoring
	Prometheus          Prometheus
	Location            string `required:"true"`
	LogLevel            string `split_words:"true" default:"INFO"`
	NumberOfMatches     int64  `split_words:"true" default:"1024"`
	ConnectionsPerMatch int64  `split_words:"true" default:"4"`
}

// Prometheus config
type Prometheus struct {
	Port int32 `default:"9102"`
}

// Monitoring config
type Monitoring struct {
	Namespace string `required:"true"`
}

// Load loads configuration from env variables
func Load() Config {
	var config Config

	if err := configlib.LoadConfig(&config, &config.App); err != nil {
		panic(err)
	}

	return config
}
