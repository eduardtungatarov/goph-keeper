package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/eduardtungatarov/goph-keeper/internal/server/repository/user/queries"
	"github.com/golang-jwt/jwt/v4"
)

// UserIDKey тип ключа в контексте хранящий id аутентифицированного пользователя.
type UserIDKey string

const (
	// UserIDKeyName наименование ключа в контексте хранящий id аутентифицированного пользователя.
	UserIDKeyName UserIDKey = "userId"
	// tokenLifeTime время жизни jwt токена.
	tokenLifeTime = 24 * time.Hour
)

// ErrLoginPwd неправильный логик и пароль.
var ErrLoginPwd = errors.New("invalid username/password pair")

// UserRepository репозиторий пользователей.
type UserRepository interface {
	SaveUser(ctx context.Context, user queries.User) (queries.User, error)
	FindUserByLogin(ctx context.Context, login string) (queries.User, error)
}

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

// Service аутентификации.
type Service struct {
	secretKey string
	userRepo  UserRepository
}

// New создать новый сервис аутентификации.
func New(secretKey string, userRepo UserRepository) *Service {
	return &Service{
		secretKey: secretKey,
		userRepo:  userRepo,
	}
}

// Register регистрация пользователя.
func (s *Service) Register(ctx context.Context, login, pwd string) (string, error) {
	hashedPassword, err := s.getHashedPwd(pwd)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.SaveUser(ctx, queries.User{
		Login:    login,
		Password: hashedPassword,
	})
	if err != nil {
		return "", err
	}

	token, err := s.buildJWTString(int(user.ID))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login аутентификация пользователя.
func (s *Service) Login(ctx context.Context, login, pwd string) (string, error) {
	user, err := s.userRepo.FindUserByLogin(ctx, login)
	if err != nil {
		return "", ErrLoginPwd
	}

	if !s.checkPasswordHash(pwd, user.Password) {
		return "", ErrLoginPwd
	}

	token, err := s.buildJWTString(int(user.ID))
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserIDByToken получить пользователя по jwt токену.
func (s *Service) GetUserIDByToken(tokenStr string) (int, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(s.secretKey), nil
		})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("auth token is not valid")
	}

	return claims.UserID, nil
}

func (s *Service) getHashedPwd(pwd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *Service) checkPasswordHash(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

func (s *Service) buildJWTString(userID int) (string, error) {
	expirationTime := time.Now().Add(tokenLifeTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
