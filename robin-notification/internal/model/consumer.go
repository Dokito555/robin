package model

type RegisterMessage struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
