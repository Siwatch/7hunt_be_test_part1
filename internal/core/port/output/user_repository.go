package output

import (
	"7hunt-be-rest-api/internal/core/domain"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	UpdateUser(ctx context.Context, userID string, user *domain.User) (int64, error)
	DeleteUser(ctx context.Context, userID string) (int64, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Count(ctx context.Context) (int64, error)
}
