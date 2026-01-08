package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/port"
)

type AuthHandler struct {
	conf *config.JWT
	svc  port.AuthService
}

func NewAuthHandler(conf *config.JWT, svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		conf,
		svc,
	}
}

type LoginUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req LoginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("email and password are required")})
		return
	}

	token, err := ah.svc.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// convert duration to int
	duration, err := strconv.Atoi(ah.conf.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// set jwt token in cookie
	c.SetCookie("jwt_token", token, duration*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
