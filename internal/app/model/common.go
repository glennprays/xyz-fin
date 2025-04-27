package model

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
