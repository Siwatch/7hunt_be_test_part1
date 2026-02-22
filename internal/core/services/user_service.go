package services

import (
	"7hunt-be-rest-api/auth"
	"7hunt-be-rest-api/internal/core/api"
	"7hunt-be-rest-api/internal/core/domain"
	"7hunt-be-rest-api/internal/core/port/input"
	"7hunt-be-rest-api/internal/core/port/output"
	"7hunt-be-rest-api/models"
	"7hunt-be-rest-api/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type userService struct {
	userRepo output.UserRepository
	hasher   utils.PasswordHasher
	logger   *utils.CustomLogger
	tokenGen auth.TokenManager
}

func NewUserService(userRepo output.UserRepository, logger *utils.CustomLogger, hasher utils.PasswordHasher, tokenGen auth.TokenManager) input.UserService {
	return &userService{
		userRepo: userRepo,
		hasher:   hasher,
		logger:   logger,
		tokenGen: tokenGen,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *api.CreateUserRequest) (*models.SuccessResponse, *models.ErrorResponse) {
	existingEmail, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		s.logger.Error("%v", err)
		return nil, models.ErrInternal
	}

	if existingEmail != nil {
		s.logger.Error("Email already exists: %v", user.Email)
		return nil, models.ErrEmailAlreadyExists
	}
	hashedPassword, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("%v", err)
		return nil, models.ErrInternal
	}

	userDomain := user.ToUserDomain(hashedPassword)
	userDomain.CreatedAt = time.Now().UTC()

	if err := s.userRepo.CreateUser(ctx, userDomain); err != nil {
		s.logger.Error("%v", err)
		return nil, models.ErrInternal
	}

	return models.Success, nil
}

func (s *userService) GetUsers(ctx context.Context) (*api.GetUsersResponse, *models.ErrorResponse) {
	res := []api.GetUserResponse{}
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		s.logger.Error("%v", err)
		return nil, models.ErrInternal
	}

	for _, user := range users {
		res = append(res, api.GetUserResponse{
			ID:        user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	return &api.GetUsersResponse{
		Users:     res,
		UserCount: len(res),
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*api.GetUserResponse, *models.ErrorResponse) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("%v", err)
		if err == domain.ErrInvalidID {
			return nil, models.ErrBadRequest
		}
		if err == domain.ErrNotFound {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrInternal
	}

	res := &api.GetUserResponse{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return res, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID string, user *api.UpdateUserRequest) (*models.SuccessResponse, *models.ErrorResponse) {
	matchedCount, err := s.userRepo.UpdateUser(ctx, userID, user.ToUserDomain())
	if err != nil {
		s.logger.Error("%v", err)
		return nil, models.ErrInternal
	}

	if matchedCount == 0 {
		s.logger.Error("%v", err)
		return nil, models.ErrUserNotFound
	}

	return models.Success, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID string) (*models.SuccessResponse, *models.ErrorResponse) {
	deletedCount, err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		s.logger.Error("%v", err)
		return nil, models.ErrInternal
	}

	if deletedCount == 0 {
		s.logger.Error("%v", err)
		return nil, models.ErrUserNotFound
	}

	return models.Success, nil
}

func (s *userService) AuthenticateUser(ctx context.Context, user *api.UserLoginRequest) (*api.UserAuthenResponse, *models.ErrorResponse) {
	currentUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		s.logger.Error("%v", err)
		return nil, models.ErrInternal
	}

	if currentUser == nil {
		return nil, models.ErrInvalidEmailOrPassword
	}

	err2 := s.hasher.ComparePassword(currentUser.Password, user.Password)
	if err2 != nil {
		s.logger.Error("%v", err2)
		return nil, models.ErrInvalidEmailOrPassword
	}

	token, err3 := s.tokenGen.GenerateToken(currentUser.ID.Hex())
	if err3 != nil {
		s.logger.Error("%v", err3)
		return nil, models.ErrInternal
	}

	return &api.UserAuthenResponse{
		Token: token,
	}, nil
}

func (s *userService) Count(ctx context.Context) (int64, *models.ErrorResponse) {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		s.logger.Error("failed to count users in repository: %v", err)
		return 0, models.ErrInternal
	}
	return count, nil
}
