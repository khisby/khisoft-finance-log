package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config struct to implement model of inbox configuration
type Config struct {
	GSheetID      string   `envconfig:"GSHEET_ID" default:""`
	WhitelistUser []string `envconfig:"WHITELIST_USER" default:""`
}

// Get to get defined configuration
func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return cfg
}
