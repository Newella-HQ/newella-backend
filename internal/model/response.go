package model

type Response struct {
	Message string `json:"message"`
	Payload any    `json:"payload"`
}
