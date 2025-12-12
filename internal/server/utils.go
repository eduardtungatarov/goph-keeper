package server

import (
	"context"
	"errors"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/auth"
)

func GetUserID(ctx context.Context) (int, error) {
	if userID, ok := ctx.Value(auth.UserIDKeyName).(int); ok {
		return userID, nil
	}
	return 0, errors.New("userID not found or not a string")
}
