package service

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/core/port"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
	notifRepo port.NotificationRepository
}

func NewUserService(repo port.UserRepository, notifRepo port.NotificationRepository) *UserService {
	return &UserService{
		repo,
		notifRepo,
	}
}

func (us *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	// hash password
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPwd

	res, err := us.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// send email notification
	if err := us.notifRepo.SendEmailNotification(ctx, res.Email); err != nil {
		return nil, err
	}

	return res, nil
}

func (us *UserService) GetUsers(ctx context.Context) ([]domain.UserResponse, error) {
	return us.repo.GetUsers(ctx)
}
