package dto

type JsonResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}
