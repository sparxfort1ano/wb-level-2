package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func newConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("LOGGER", &cfg); err != nil {
		return config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

// NewConfigMust builds the logger configuration.
// If there are errors, it panics. Panic is allowed:
// the logger middleware of the application cannot function
// without a running logger.
func NewConfigMust() config {
	config, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	return config
}
