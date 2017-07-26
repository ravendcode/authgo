package main

import "os"

// Config App
type Config struct {
	Env       string
	Port      string
	StaticDir string
	AppName   string
}

// NewConfig create new instance Config struct
func NewConfig() *Config {
	config := new(Config)
	config.AppName = "Auth"
	config.Port = ":80"
	config.StaticDir = "./static"
	env := os.Getenv("ENV")
	if env != "" {
		config.Env = env
	} else {
		config.Env = "production"
	}

	return config
}
