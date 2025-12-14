package config

import "os"

type Config struct {
	DatabaseDSN    string
	GRPCPort       string
	JWTSecretKey   string
	DataEncryptKey string
	ServerCert     string
	ServerKey      string
}

func Load() *Config {
	dsn, _ := os.LookupEnv("DATABASE_DSN")
	grpcPort, _ := os.LookupEnv("GRPC_PORT")
	JWTSecretKey, _ := os.LookupEnv("JWT_SECRET_KEY")
	dataEncryptKey, _ := os.LookupEnv("DATA_ENCRYPT_KEY")
	servertCert, _ := os.LookupEnv("SERVER_CERT")
	serverKey, _ := os.LookupEnv("SERVER_KEY")

	return &Config{
		DatabaseDSN:    dsn,
		GRPCPort:       grpcPort,
		JWTSecretKey:   JWTSecretKey,
		DataEncryptKey: dataEncryptKey,
		ServerKey:      serverKey,
		ServerCert:     servertCert,
	}
}
