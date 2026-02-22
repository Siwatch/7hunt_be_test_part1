package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"7hunt-be-rest-api/internal/core/api"
	"7hunt-be-rest-api/models"
	"7hunt-be-rest-api/validator"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) AuthenticateUser(ctx context.Context, user *api.UserLoginRequest) (*api.UserAuthenResponse, *models.ErrorResponse) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*models.ErrorResponse)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*api.UserAuthenResponse), nil
	}
	return args.Get(0).(*api.UserAuthenResponse), args.Get(1).(*models.ErrorResponse)
}

func (m *MockUserService) CreateUser(ctx context.Context, user *api.CreateUserRequest) (*models.SuccessResponse, *models.ErrorResponse) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*models.ErrorResponse)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*models.SuccessResponse), nil
	}
	return args.Get(0).(*models.SuccessResponse), args.Get(1).(*models.ErrorResponse)
}

func (m *MockUserService) GetUsers(ctx context.Context) (*api.GetUsersResponse, *models.ErrorResponse) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*models.ErrorResponse)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*api.GetUsersResponse), nil
	}
	return args.Get(0).(*api.GetUsersResponse), args.Get(1).(*models.ErrorResponse)
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID string) (*api.GetUserResponse, *models.ErrorResponse) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*models.ErrorResponse)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*api.GetUserResponse), nil
	}
	return args.Get(0).(*api.GetUserResponse), args.Get(1).(*models.ErrorResponse)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID string, user *api.UpdateUserRequest) (*models.SuccessResponse, *models.ErrorResponse) {
	args := m.Called(ctx, userID, user)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*models.ErrorResponse)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*models.SuccessResponse), nil
	}
	return args.Get(0).(*models.SuccessResponse), args.Get(1).(*models.ErrorResponse)
}

func (m *MockUserService) DeleteUser(ctx context.Context, userID string) (*models.SuccessResponse, *models.ErrorResponse) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*models.ErrorResponse)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*models.SuccessResponse), nil
	}
	return args.Get(0).(*models.SuccessResponse), args.Get(1).(*models.ErrorResponse)
}

func (m *MockUserService) Count(ctx context.Context) (int64, *models.ErrorResponse) {
	args := m.Called(ctx)
	if args.Get(1) == nil {
		return args.Get(0).(int64), nil
	}
	return args.Get(0).(int64), args.Get(1).(*models.ErrorResponse)
}

func TestUserHandler_RegisterUser(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	t.Run("UserHandler_RegisterUser_Success", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		reqBody := api.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockService.On("CreateUser", mock.Anything, &reqBody).Return(models.Success, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterUser(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.Success.Message, response.Message)
		mockService.AssertExpectations(t)
	})

	t.Run("UserHandler_RegisterUser_BindJSON_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON format
		invalidJSON := `{"name": "Test User", "email": "test@example.com", "password": "password123"`
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrBadRequest.StatusCode, response.StatusCode)
		assert.Equal(t, models.ErrBadRequest.Error.ErrorDesc, response.Error.ErrorDesc)
		mockService.AssertExpectations(t)
	})

	t.Run("UserHandler_RegisterUser_Validation_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		// Missing required field "Password"
		reqBody := api.CreateUserRequest{
			Name:  "Test User",
			Email: "test@example.com",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterUser(c)

		// handler expects StatusBadRequest for validation error
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, models.ErrValidateMessage, response.Error.ErrorDesc)
		assert.NotNil(t, response.Error.ErrorValidate)
		mockService.AssertExpectations(t)
	})

	t.Run("UserHandler_RegisterUser_Service_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		reqBody := api.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		// Service returns an error e.g. Email Already Exists
		mockService.On("CreateUser", mock.Anything, &reqBody).Return(nil, models.ErrEmailAlreadyExists)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterUser(c)

		// Handler responds with the error's status code
		assert.Equal(t, models.ErrEmailAlreadyExists.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrEmailAlreadyExists.StatusCode, response.StatusCode)
		assert.Equal(t, models.ErrEmailAlreadyExists.Error.ErrorDesc, response.Error.ErrorDesc)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("UserHandler_GetUsers_Success", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		mockResponse := &api.GetUsersResponse{
			Users: []api.GetUserResponse{
				{ID: "1", Name: "User 1", Email: "user1@example.com"},
			},
			UserCount: 1,
		}

		mockService.On("GetUsers", mock.Anything).Return(mockResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/users", nil)

		handler.GetUsers(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response api.GetUsersResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.UserCount)
		assert.Equal(t, "User 1", response.Users[0].Name)
		mockService.AssertExpectations(t)
	})

	t.Run("UserHandler_GetUsers_Service_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		mockService.On("GetUsers", mock.Anything).Return(nil, models.ErrInternal)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/users", nil)

		handler.GetUsers(c)

		assert.Equal(t, models.ErrInternal.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrInternal.StatusCode, response.StatusCode)
		assert.Equal(t, models.ErrInternal.Error.ErrorDesc, response.Error.ErrorDesc)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("UserHandler_GetUserByID_Success", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		userID := "123"
		mockResponse := &api.GetUserResponse{
			ID:    userID,
			Name:  "Test User",
			Email: "test@example.com",
		}

		mockService.On("GetUserByID", mock.Anything, userID).Return(mockResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		c.Params = []gin.Param{{Key: "userId", Value: userID}}

		handler.GetUserByID(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response api.GetUserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, userID, response.ID)
		assert.Equal(t, "Test User", response.Name)
		mockService.AssertExpectations(t)
	})

	t.Run("UserHandler_GetUserByID_ParamRequired_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/users/", nil)
		c.Params = []gin.Param{{Key: "userId", Value: ""}} // Missing ID

		handler.GetUserByID(c)

		// Expect custom ErrParamIdRequired code (400)
		assert.Equal(t, models.ErrParamIdRequired.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrParamIdRequired.StatusCode, response.StatusCode)
		assert.Equal(t, models.ErrParamIdRequired.Error.ErrorDesc, response.Error.ErrorDesc)
	})

	t.Run("UserHandler_GetUserByID_Service_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		userID := "123"
		mockService.On("GetUserByID", mock.Anything, userID).Return(nil, models.ErrUserNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		c.Params = []gin.Param{{Key: "userId", Value: userID}}

		handler.GetUserByID(c)

		assert.Equal(t, models.ErrUserNotFound.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrUserNotFound.StatusCode, response.StatusCode)
		assert.Equal(t, models.ErrUserNotFound.Error.ErrorDesc, response.Error.ErrorDesc)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("UserHandler_UpdateUser_Success", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		userID := "123"
		reqBody := api.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		mockService.On("UpdateUser", mock.Anything, userID, &reqBody).Return(models.Success, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/users/"+userID, bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "userId", Value: userID}}

		handler.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.Success.Message, response.Message)
		mockService.AssertExpectations(t)
	})

	t.Run("UserHandler_UpdateUser_ParamRequired_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPut, "/api/users/", nil)
		c.Params = []gin.Param{{Key: "userId", Value: ""}}

		handler.UpdateUser(c)

		assert.Equal(t, models.ErrParamIdRequired.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrParamIdRequired.StatusCode, response.StatusCode)
	})

	t.Run("UserHandler_UpdateUser_BindJSON_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		invalidJSON := `{"name": "Updated Name"`
		c.Request = httptest.NewRequest(http.MethodPut, "/api/users/123", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "userId", Value: "123"}}

		handler.UpdateUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrBadRequest.StatusCode, response.StatusCode)
	})

	t.Run("UserHandler_UpdateUser_Service_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		userID := "123"
		reqBody := api.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		mockService.On("UpdateUser", mock.Anything, userID, &reqBody).Return(nil, models.ErrUserNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/api/users/"+userID, bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "userId", Value: userID}}

		handler.UpdateUser(c)

		assert.Equal(t, models.ErrUserNotFound.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrUserNotFound.StatusCode, response.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("UserHandler_DeleteUser_Success", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		userID := "123"
		mockService.On("DeleteUser", mock.Anything, userID).Return(models.Success, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/users/"+userID, nil)
		c.Params = []gin.Param{{Key: "userId", Value: userID}}

		handler.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.Success.Message, response.Message)
		mockService.AssertExpectations(t)
	})

	t.Run("UserHandler_DeleteUser_ParamRequired_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/users/", nil)
		c.Params = []gin.Param{{Key: "userId", Value: ""}}

		handler.DeleteUser(c)

		assert.Equal(t, models.ErrParamIdRequired.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrParamIdRequired.StatusCode, response.StatusCode)
	})

	t.Run("UserHandler_DeleteUser_Service_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		userID := "123"
		mockService.On("DeleteUser", mock.Anything, userID).Return(nil, models.ErrUserNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/users/"+userID, nil)
		c.Params = []gin.Param{{Key: "userId", Value: userID}}

		handler.DeleteUser(c)

		assert.Equal(t, models.ErrUserNotFound.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrUserNotFound.StatusCode, response.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_LoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("UserHandler_LoginUser_Success", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		reqBody := api.UserLoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockResponse := &api.UserAuthenResponse{
			Token: "test_jwt_token",
		}

		mockService.On("AuthenticateUser", mock.Anything, &reqBody).Return(mockResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginUser(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response api.UserAuthenResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test_jwt_token", response.Token)
		mockService.AssertExpectations(t)
	})

	t.Run("UserHandler_LoginUser_BindJSON_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		invalidJSON := `{"email": "test@example.com", "password": "password123"`
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrBadRequest.StatusCode, response.StatusCode)
	})

	t.Run("UserHandler_LoginUser_Validation_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		reqBody := api.UserLoginRequest{
			Email: "test@example.com",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, models.ErrValidateMessage, response.Error.ErrorDesc)
		assert.NotNil(t, response.Error.ErrorValidate)
	})

	t.Run("UserHandler_LoginUser_Service_Error", func(t *testing.T) {
		mockService := new(MockUserService)
		val := validator.NewValidator()
		handler := NewUserHandler(mockService, val)

		reqBody := api.UserLoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockService.On("AuthenticateUser", mock.Anything, &reqBody).Return(nil, models.ErrInvalidEmailOrPassword)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginUser(c)

		assert.Equal(t, models.ErrInvalidEmailOrPassword.StatusCode, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, models.ErrInvalidEmailOrPassword.StatusCode, response.StatusCode)
		mockService.AssertExpectations(t)
	})
}
