package token

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Storage struct {
	Token string `json:"token"`
}

func Save(token string) error {
	dir := filepath.Join(os.Getenv("HOME"), ".keeper")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	file := filepath.Join(dir, "token.json")
	data := Storage{Token: token}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, bytes, 0o600)
}

func Load() (string, error) {
	file := filepath.Join(os.Getenv("HOME"), ".keeper", "token.json")
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	var ts Storage
	if err := json.Unmarshal(data, &ts); err != nil {
		return "", err
	}
	return ts.Token, nil
}

func Clear() error {
	file := filepath.Join(os.Getenv("HOME"), ".keeper", "token.json")
	return os.Remove(file)
}
