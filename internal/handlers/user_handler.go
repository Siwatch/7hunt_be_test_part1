package handlers

import (
	"7hunt-be-rest-api/internal/core/api"
	"7hunt-be-rest-api/internal/core/port/input"
	"7hunt-be-rest-api/models"
	"7hunt-be-rest-api/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService input.UserService
	validator   *validator.Validator
}

func NewUserHandler(userService input.UserService, validator *validator.Validator) *UserHandler {
	return &UserHandler{
		validator:   validator,
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	req := new(api.CreateUserRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrBadRequest)
		return
	}

	if err := h.validator.ValidateRequest(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrValidate(err))
		return
	}

	res, errRes := h.userService.CreateUser(ctx, req)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	res, errRes := h.userService.GetUsers(ctx)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(models.ErrParamIdRequired.StatusCode, models.ErrParamIdRequired)
		return
	}

	res, errRes := h.userService.GetUserByID(ctx, userID)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(models.ErrParamIdRequired.StatusCode, models.ErrParamIdRequired)
		return
	}

	req := new(api.UpdateUserRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrBadRequest)
		return
	}

	if err := h.validator.ValidateRequest(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrValidate(err))
		return
	}

	res, errRes := h.userService.UpdateUser(ctx, userID, req)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(models.ErrParamIdRequired.StatusCode, models.ErrParamIdRequired)
		return
	}

	res, errRes := h.userService.DeleteUser(ctx, userID)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) LoginUser(ctx *gin.Context) {
	req := new(api.UserLoginRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrBadRequest)
		return
	}

	if err := h.validator.ValidateRequest(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrValidate(err))
		return
	}

	res, errRes := h.userService.AuthenticateUser(ctx, req)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
