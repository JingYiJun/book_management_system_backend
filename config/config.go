package config

import (
	"github.com/caarlos0/env/v6"
	"net/url"
)

var Config struct {
	Mode        string  `env:"MODE" envDefault:"dev"`
	Debug       bool    `env:"DEBUG" envDefault:"false"`
	PostgresDSN url.URL `env:"POSTGRES_DSN"`
	AppName     string  `env:"APP_NAME" envDefault:"book_management_system"`
	Hostname    string  `env:"HOSTNAME" envDefault:"localhost"`
}

func InitConfig() {
	var err error
	if err = env.Parse(&Config); err != nil {
		panic(err)
	}
}

const (
	ModeTest       = "test"
	ModeDev        = "dev"
	ModeProduction = "production"
	ModeBench      = "bench"
)
