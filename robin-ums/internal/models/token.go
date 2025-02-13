package model

import "github.com/golang-jwt/jwt/v4"

type ClaimToken struct {
	UserID 	int 	`json:"user_id"`
	Email 	string	`json:"email"`
	Role	string	`json:"role"`
	jwt.RegisteredClaims
}