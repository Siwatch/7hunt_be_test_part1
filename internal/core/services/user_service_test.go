package services

import (
	"7hunt-be-rest-api/auth"
	"7hunt-be-rest-api/internal/core/api"
	"7hunt-be-rest-api/internal/core/domain"
	"7hunt-be-rest-api/models"
	"7hunt-be-rest-api/utils"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockUserRepository struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) ValidateToken(tokenString string) (*auth.Claims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*auth.Claims), args.Error(1)
}

func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) ComparePassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, userID string, user *domain.User) (int64, error) {
	args := m.Called(ctx, userID, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("UserService_CreateUser_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)
		mockPasswordHasher.On("HashPassword", req.Password).Return("hashed_password", nil)
		mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

		res, errRes := service.CreateUser(context.Background(), req)

		assert.Nil(t, errRes)
		assert.Equal(t, models.Success, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_CreateUser_Internal_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, errors.New("Error Internal"))

		res, errRes := service.CreateUser(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_CreateUser_Email_Already_Exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.CreateUserRequest{
			Name:     "Test User",
			Email:    "existing@example.com",
			Password: "password123",
		}

		existingUser := &domain.User{Email: req.Email}
		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

		res, errRes := service.CreateUser(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrEmailAlreadyExists, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_CreateUser_Password_Hashing_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)
		mockPasswordHasher.On("HashPassword", req.Password).Return("", errors.New("hashing error"))

		res, errRes := service.CreateUser(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_CreateUser_Database_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)
		mockPasswordHasher.On("HashPassword", req.Password).Return("hashed_password", nil)
		mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("database error"))

		res, errRes := service.CreateUser(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUsers(t *testing.T) {
	t.Run("UserService_GetUsers_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		users := []*domain.User{
			{Name: "Alice", Email: "alice@example.com"},
			{Name: "Bob", Email: "bob@example.com"},
		}

		mockRepo.On("GetUsers", mock.Anything).Return(users, nil)

		res, errRes := service.GetUsers(context.Background())

		assert.Nil(t, errRes)
		assert.NotNil(t, res)
		assert.Equal(t, 2, res.UserCount)
		assert.Equal(t, "Alice", res.Users[0].Name)
		assert.Equal(t, "Bob", res.Users[1].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_GetUsers_Database_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		mockRepo.On("GetUsers", mock.Anything).Return([]*domain.User{}, errors.New("database error"))

		res, errRes := service.GetUsers(context.Background())

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	t.Run("UserService_GetUserByID_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "507f1f77bcf86cd799439011"
		user := &domain.User{Name: "Alice", Email: "alice@example.com"}

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)

		res, errRes := service.GetUserByID(context.Background(), userID)

		assert.Nil(t, errRes)
		assert.NotNil(t, res)
		assert.Equal(t, "Alice", res.Name)
		assert.Equal(t, "alice@example.com", res.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_GetUserByID_Invalid_ID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "invalid-id"

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(nil, domain.ErrInvalidID)

		res, errRes := service.GetUserByID(context.Background(), userID)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrBadRequest, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_GetUserByID_Not_Found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "507f1f77bcf86cd799439011"

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(nil, domain.ErrNotFound)

		res, errRes := service.GetUserByID(context.Background(), userID)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrUserNotFound, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_GetUserByID_Database_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "507f1f77bcf86cd799439011"

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(nil, errors.New("database error"))

		res, errRes := service.GetUserByID(context.Background(), userID)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Run("UserService_UpdateUser_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "abc123"
		req := &api.UpdateUserRequest{Name: "Updated Name", Email: "updated@example.com"}

		mockRepo.On("UpdateUser", mock.Anything, userID, mock.AnythingOfType("*domain.User")).Return(int64(1), nil)

		res, errRes := service.UpdateUser(context.Background(), userID, req)

		assert.Nil(t, errRes)
		assert.Equal(t, models.Success, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_UpdateUser_User_Not_Found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "notexist"
		req := &api.UpdateUserRequest{Name: "Updated Name", Email: "updated@example.com"}

		mockRepo.On("UpdateUser", mock.Anything, userID, mock.AnythingOfType("*domain.User")).Return(int64(0), nil)

		res, errRes := service.UpdateUser(context.Background(), userID, req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrUserNotFound, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_UpdateUser_Database_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "abc123"
		req := &api.UpdateUserRequest{Name: "Updated Name", Email: "updated@example.com"}

		mockRepo.On("UpdateUser", mock.Anything, userID, mock.AnythingOfType("*domain.User")).Return(int64(0), errors.New("database error"))

		res, errRes := service.UpdateUser(context.Background(), userID, req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("UserService_DeleteUser_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "abc123"

		mockRepo.On("DeleteUser", mock.Anything, userID).Return(int64(1), nil)

		res, errRes := service.DeleteUser(context.Background(), userID)

		assert.Nil(t, errRes)
		assert.Equal(t, models.Success, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_DeleteUser_User_Not_Found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "notexist"

		mockRepo.On("DeleteUser", mock.Anything, userID).Return(int64(0), nil)

		res, errRes := service.DeleteUser(context.Background(), userID)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrUserNotFound, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_DeleteUser_Database_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		userID := "abc123"

		mockRepo.On("DeleteUser", mock.Anything, userID).Return(int64(0), errors.New("database error"))

		res, errRes := service.DeleteUser(context.Background(), userID)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_AuthenticateUser(t *testing.T) {
	t.Run("UserService_AuthenticateUser_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.UserLoginRequest{Email: "test@example.com", Password: "password123"}
		existingUser := &domain.User{Email: req.Email, Password: "hashed_password"}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)
		mockPasswordHasher.On("ComparePassword", existingUser.Password, req.Password).Return(nil)
		mockTokenManager.On("GenerateToken", existingUser.ID.Hex()).Return("token_string", nil)

		res, errRes := service.AuthenticateUser(context.Background(), req)

		assert.Nil(t, errRes)
		assert.NotNil(t, res)
		assert.Equal(t, "token_string", res.Token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_AuthenticateUser_User_Not_Found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.UserLoginRequest{Email: "notfound@example.com", Password: "password123"}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)

		res, errRes := service.AuthenticateUser(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInvalidEmailOrPassword, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_AuthenticateUser_Internal_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.UserLoginRequest{Email: "test@example.com", Password: "password123"}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, errors.New("database error"))

		res, errRes := service.AuthenticateUser(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_AuthenticateUser_Invalid_Password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.UserLoginRequest{Email: "test@example.com", Password: "wrongpassword"}
		existingUser := &domain.User{Email: req.Email, Password: "hashed_password"}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)
		mockPasswordHasher.On("ComparePassword", existingUser.Password, req.Password).Return(errors.New("invalid password"))

		res, errRes := service.AuthenticateUser(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInvalidEmailOrPassword, errRes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_AuthenticateUser_Generate_Token_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		req := &api.UserLoginRequest{Email: "test@example.com", Password: "password123"}
		existingUser := &domain.User{Email: req.Email, Password: "hashed_password"}

		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)
		mockPasswordHasher.On("ComparePassword", existingUser.Password, req.Password).Return(nil)
		mockTokenManager.On("GenerateToken", existingUser.ID.Hex()).Return("", errors.New("token error"))

		res, errRes := service.AuthenticateUser(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Count(t *testing.T) {
	t.Run("UserService_Count_Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		mockRepo.On("Count", mock.Anything).Return(int64(5), nil)

		count, errRes := service.Count(context.Background())

		assert.Nil(t, errRes)
		assert.Equal(t, int64(5), count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserService_Count_Database_Error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenManager := new(MockTokenManager)
		mockPasswordHasher := new(MockPasswordHasher)
		logger := utils.NewLogger("user-service")
		service := NewUserService(mockRepo, logger, mockPasswordHasher, mockTokenManager)

		mockRepo.On("Count", mock.Anything).Return(int64(0), errors.New("database error"))

		count, errRes := service.Count(context.Background())

		assert.Equal(t, int64(0), count)
		assert.Equal(t, models.ErrInternal, errRes)
		mockRepo.AssertExpectations(t)
	})
}
