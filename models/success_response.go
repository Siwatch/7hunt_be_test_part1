package models

type SuccessResponse struct {
	Message string `json:"Message"`
}

var (
	Success = &SuccessResponse{Message: "Success"}
)
