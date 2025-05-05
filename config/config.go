package config

import (
	"github.com/Kry0z1/impulse/eventtime"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Laps        int                   `json:"laps" env-required:"true"`
	LapLen      int                   `json:"lapLen" env-required:"true"`
	PenaltyLen  int                   `json:"penaltyLen" env-required:"true"`
	FiringLines int                   `json:"firingLines" env-required:"true"`
	Start       eventtime.TimestampMS `json:"start" env-required:"true"`
	StartDelta  eventtime.Timestamp   `json:"startDelta" env-required:"true"`
}

func MustLoad(path string) *Config {
	if path == "" {
		panic("empty config path")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("couldn't read config: " + err.Error())
	}

	return &cfg
}
