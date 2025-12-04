package config

import "os"

type Config struct {
	DatabaseDSN string
	GRPCPort    string
}

func Load() *Config {
	dsn, _ := os.LookupEnv("DATABASE_DSN")
	grpcPort, _ := os.LookupEnv("GRPC_PORT")

	return &Config{
		DatabaseDSN: dsn,
		GRPCPort:    grpcPort,
	}
}
