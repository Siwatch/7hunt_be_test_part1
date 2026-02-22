package models

// ใส่ * ทำให้ return  nil ได้ถ้าไม่มี error
type ErrorResponse struct {
	StatusCode int                `json:"StatusCode"`
	Error      *ErrorResponseBody `json:"Error,omitempty"`
}

type ErrorResponseBody struct {
	ErrorDesc     string  `json:"ErrorDesc"`
	ErrorValidate *string `json:"ErrorValidate,omitempty"`
}

type ErrorValidate struct {
	Field string `json:"Field"`
	Desc  string `json:"Desc"`
}
