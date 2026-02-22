package input

import (
	"7hunt-be-rest-api/internal/core/api"
	"7hunt-be-rest-api/models"
	"context"
)

type UserService interface {
	AuthenticateUser(ctx context.Context, user *api.UserLoginRequest) (*api.UserAuthenResponse, *models.ErrorResponse)
	CreateUser(ctx context.Context, user *api.CreateUserRequest) (*models.SuccessResponse, *models.ErrorResponse)
	GetUsers(ctx context.Context) (*api.GetUsersResponse, *models.ErrorResponse)
	GetUserByID(ctx context.Context, userID string) (*api.GetUserResponse, *models.ErrorResponse)
	UpdateUser(ctx context.Context, userID string, user *api.UpdateUserRequest) (*models.SuccessResponse, *models.ErrorResponse)
	DeleteUser(ctx context.Context, userID string) (*models.SuccessResponse, *models.ErrorResponse)
	Count(ctx context.Context) (int64, *models.ErrorResponse)
}
