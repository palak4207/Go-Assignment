package db

import (
	"context"
)

// Store ...
type Store interface {
	UserGetter
	UserUpdater
	UserCreator
}

type UserGetter interface {
	GetUserByEmail(ctx context.Context, email string) (UserResponse, error)
	GetUserByUserName(ctx context.Context, userName string) (UserResponse, error)
	GetUsers(ctx context.Context) ([]UserResponse, error)
}

type UserCreator interface {
	CreateUser(ctx context.Context, user UserRequest) (UserResponse, error)
}

type UserUpdater interface {
	UpdateUser(ctx context.Context, user UserRequest) (UserResponse, error)
}
