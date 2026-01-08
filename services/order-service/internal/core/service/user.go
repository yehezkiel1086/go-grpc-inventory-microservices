package service

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/port"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo,
	}
}

func (us *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPwd

	return us.repo.CreateUser(ctx, user)
}

func (us *UserService) GetUsers(ctx context.Context) ([]domain.UserResponse, error) {
	return us.repo.GetUsers(ctx)
}
