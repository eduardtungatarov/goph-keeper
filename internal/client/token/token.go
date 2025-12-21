package token

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// TokenStorage хранилище токена после входа.
type TokenStorage struct{}

type tokenData struct {
	Token string `json:"token"`
}

func New() *TokenStorage {
	return &TokenStorage{}
}

// Save сохранить токен в файл в домашнюю директорию пользователя.
func (t *TokenStorage) Save(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(homeDir, ".keeper")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	file := filepath.Join(dir, "token.json")
	data := tokenData{Token: token}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, bytes, 0o600)
}

// Load загрузить токен в приложение из файла.
func (t *TokenStorage) Load() (string, error) {
	file := filepath.Join(os.Getenv("HOME"), ".keeper", "token.json")
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	var ts tokenData
	if err := json.Unmarshal(data, &ts); err != nil {
		return "", err
	}
	return ts.Token, nil
}
