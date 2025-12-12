package config

import "os"

type Config struct {
	DatabaseDSN    string
	GRPCPort       string
	JWTSecretKey   string
	DataEncryptKey string
}

func Load() *Config {
	dsn, _ := os.LookupEnv("DATABASE_DSN")
	grpcPort, _ := os.LookupEnv("GRPC_PORT")
	JWTSecretKey, _ := os.LookupEnv("JWT_SECRET_KEY")
	DataEncryptKey, _ := os.LookupEnv("DATA_ENCRYPT_KEY")

	return &Config{
		DatabaseDSN:    dsn,
		GRPCPort:       grpcPort,
		JWTSecretKey:   JWTSecretKey,
		DataEncryptKey: DataEncryptKey,
	}
}
