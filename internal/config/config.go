package config

import "os"

type Config struct {
	DatabaseDSN  string
	GRPCPort     string
	JWTSecretKey string
}

func Load() *Config {
	dsn, _ := os.LookupEnv("DATABASE_DSN")
	grpcPort, _ := os.LookupEnv("GRPC_PORT")
	JWTSecretKey, _ := os.LookupEnv("JWT_SECRET_KEY")

	return &Config{
		DatabaseDSN:  dsn,
		GRPCPort:     grpcPort,
		JWTSecretKey: JWTSecretKey,
	}
}
