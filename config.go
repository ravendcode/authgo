package main

import (
	"os"
)

// Config App
type Config struct {
	Env       string
	Port      string
	StaticDir string
	AppName   string
	Secret    string
	RemeberMe int
}

// NewConfig create new instance Config struct
func NewConfig() *Config {
	config := new(Config)
	config.AppName = "Auth"
	config.Port = ":80"
	config.StaticDir = "./static"
	// Login on 1 year
	config.RemeberMe = 86400 * 365

	if env := os.Getenv("ENV"); env != "" {
		config.Env = env
	} else {
		config.Env = "production"
	}

	if secret := os.Getenv("SECRET"); secret != "" {
		config.Secret = secret
	} else {
		config.Secret = "472c8bdb5fac3b79a938a02e1879bbd8c687900875e6bafef44e9ef4401fc519"
	}

	return config
}
