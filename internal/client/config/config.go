package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	ServerCert    string
	ServerTimeout time.Duration
	ServerAddress string
}

func Load() (*Config, error) {
	servertCert, ok := os.LookupEnv("SERVER_CERT")
	if !ok {
		return nil, errors.New("SERVER_CERT not defined")
	}
	serverAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		return nil, errors.New("SERVER_ADDRESS not defined")
	}

	return &Config{
		ServerCert:    servertCert,
		ServerTimeout: time.Second * 5,
		ServerAddress: serverAddress,
	}, nil
}
