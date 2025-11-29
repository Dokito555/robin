package model

type BaseResponse[T any] struct {
	Message int    `json:"message"`
	Error   string `json:"error,omitempty"`
	Data    T      `json:"data"`
}