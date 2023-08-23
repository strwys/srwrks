package model

type APIResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseError struct {
	Message string `json:"errors"`
}
