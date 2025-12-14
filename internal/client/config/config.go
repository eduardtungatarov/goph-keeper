package config

import (
	"os"
	"time"
)

type Config struct {
	ServerCert    string
	ServerTimeout time.Duration
	ServerAddress string
}

func Load() *Config {
	servertCert, _ := os.LookupEnv("SERVER_CERT")
	serverAddress, _ := os.LookupEnv("SERVER_ADDRESS")

	return &Config{
		ServerCert:    servertCert,
		ServerTimeout: time.Second * 5,
		ServerAddress: serverAddress,
	}
}
