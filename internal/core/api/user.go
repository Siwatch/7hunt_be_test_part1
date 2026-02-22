package api

import (
	"7hunt-be-rest-api/internal/core/domain"
	"time"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GetUserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUsersResponse struct {
	Users     []GetUserResponse `json:"users"`
	UserCount int               `json:"user_count"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserAuthenResponse struct {
	Token string `json:"Token"`
}

func (r *CreateUserRequest) ToUserDomain(hashedPassword string) *domain.User {
	return &domain.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: hashedPassword,
	}
}

func (r *UpdateUserRequest) ToUserDomain() *domain.User {
	return &domain.User{
		Name:  r.Name,
		Email: r.Email,
	}
}
