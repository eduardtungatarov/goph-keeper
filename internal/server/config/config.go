package config

import "os"

// Config сервера.
type Config struct {
	DatabaseDSN    string // строка подключения к базе данных
	ServerAddress  string // хост:порт запуска сервера
	JWTSecretKey   string // секретный ключ для jwt-токенов
	DataEncryptKey string // секретный ключ для шифрования данных пользователей
	ServerCert     string // публичный сертификат сервера
	ServerKey      string // приватный ключ сервера
}

// Load загрузка конфига из env.
func Load() *Config {
	dsn, _ := os.LookupEnv("DATABASE_DSN")
	serverAddress, _ := os.LookupEnv("SERVER_ADDRESS")
	JWTSecretKey, _ := os.LookupEnv("JWT_SECRET_KEY")
	dataEncryptKey, _ := os.LookupEnv("DATA_ENCRYPT_KEY")
	servertCert, _ := os.LookupEnv("SERVER_CERT")
	serverKey, _ := os.LookupEnv("SERVER_KEY")

	return &Config{
		DatabaseDSN:    dsn,
		ServerAddress:  serverAddress,
		JWTSecretKey:   JWTSecretKey,
		DataEncryptKey: dataEncryptKey,
		ServerKey:      serverKey,
		ServerCert:     servertCert,
	}
}
