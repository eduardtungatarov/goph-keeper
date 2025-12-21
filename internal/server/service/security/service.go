package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// Service сервис шифрования данных пользователя.
type Service struct {
	key []byte
}

// New конструктор сервиса шифрования данных пользователя.
func New(key []byte) *Service {
	return &Service{
		key: key,
	}
}

// GetEncrypted получить зашифрованную форму данных.
func (s *Service) GetEncrypted(ctx context.Context, data []byte) ([]byte, error) {
	const op = "security.Service.GetEncrypted"

	aesblock, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	nonce, err := s.generateRandom(aesgcm.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	encrypted := aesgcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

// GetDecrypted дешифровать данные пользователя.
func (s *Service) GetDecrypted(ctx context.Context, data []byte) ([]byte, error) {
	const op = "security.Service.GetDecrypted"

	aesblock, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("%s: data too short", op)
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	res, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *Service) generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
