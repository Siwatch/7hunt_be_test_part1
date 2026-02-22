package models

const (
	ErrValidateMessage               = "Validate Failed"
	ErrInternalMessage               = "Internal Server Error"
	ErrBadRequestMessage             = "Bad Request"
	ErrParamIdRequiredMessage        = "Param userId is required"
	ErrUserNotFoundMessage           = "User not found"
	ErrEmailAlreadyExistsMessage     = "Email already exists"
	ErrInvalidPasswordMessage        = "Invalid password"
	ErrInvalidEmailOrPasswordMessage = "Invalid email or password"
	ErrUnAuthorizeMessage            = "UnAuthorize"
	ErrInvalidTokenMessage           = "Invalid Token"
	ErrTimeoutMessage                = "Timeout"
)

var (
	ErrInternal               = NewError(500, ErrInternalMessage)
	ErrBadRequest             = NewError(400, ErrBadRequestMessage)
	ErrParamIdRequired        = NewError(400, ErrParamIdRequiredMessage)
	ErrUserNotFound           = NewError(404, ErrUserNotFoundMessage)
	ErrEmailAlreadyExists     = NewError(400, ErrEmailAlreadyExistsMessage)
	ErrInvalidPassword        = NewError(400, ErrInvalidPasswordMessage)
	ErrInvalidEmailOrPassword = NewError(400, ErrInvalidEmailOrPasswordMessage)
	ErrUnAuthorize            = NewError(401, ErrUnAuthorizeMessage)
	ErrInvalidToken           = NewError(401, ErrInvalidTokenMessage)
	ErrTimeout                = NewError(408, ErrTimeoutMessage)
)

func NewError(statusCode int, desc string) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: statusCode,
		Error: &ErrorResponseBody{
			ErrorDesc: desc,
		},
	}
}

func ErrValidate(errValidate *string) *ErrorResponse {
	err := NewError(400, ErrValidateMessage)
	err.Error.ErrorValidate = errValidate
	return err
}
