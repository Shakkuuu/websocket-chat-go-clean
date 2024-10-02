package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	UserNameDB string `env:"DB_USERNAME"`
	UserPassDB string `env:"DB_USERPASS"`
	ProtocolDB string `env:"DB_PROTOCOL"`
	NameDB     string `env:"DB_DATABASENAME"`
	PortDB     string `env:"DB_PORT"`
	Port       string `env:"SERVERPORT"`
	SessionKey string `env:"SESSION_KEY"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
