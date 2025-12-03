package config

import "os"

type Config struct {
	DatabaseDSN string
}

func Load() Config {
	dsn, _ := os.LookupEnv("DATABASE_DSN")

	return Config{
		DatabaseDSN: dsn,
	}
}
