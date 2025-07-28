package context

import (
	"context"
	"lenslocked.com/models"
)

type userCtxKeyType string

const userCtxKey userCtxKeyType = "user"

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userCtxKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
