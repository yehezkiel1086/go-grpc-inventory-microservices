package util

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/core/domain"
)

func GenerateJWTToken(conf *config.JWT, user *domain.User) (string, error) {
	mySigningKey := []byte(conf.Secret)

	// convert duration to int
	duration, err := strconv.Atoi(conf.Duration)
	if err != nil {
		return "", err
	}

	// Create claims with multiple fields populated
	claims := domain.JWTClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(duration) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}
