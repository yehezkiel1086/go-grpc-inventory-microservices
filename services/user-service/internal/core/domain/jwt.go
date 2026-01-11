package domain

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
	jwt.RegisteredClaims
}
