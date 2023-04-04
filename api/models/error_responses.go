package models

type BadRequestResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

type NotFoundResponse struct {
	Error string `json:"error"`
}

type InternalServerErrorResponse struct {
	Error string `json:"error"`
}
