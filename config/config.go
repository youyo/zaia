package config

import "time"

const (
	RequestTimeout = 10 * time.Second
	TimeLayout     = "2006-01-02T15:04:05Z"
)

type Config struct {
	RequestTimeout time.Duration
	TimeLayout     string
}

func NewConfig() *Config {
	return &Config{
		RequestTimeout: RequestTimeout,
		TimeLayout:     TimeLayout,
	}
}
