package repository

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	db := ur.db.GetDB()
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	db := ur.db.GetDB()

	var user *domain.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUsers(ctx context.Context) ([]domain.UserResponse, error) {
	db := ur.db.GetDB()

	var users []domain.UserResponse
	if err := db.Model(&domain.User{}).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
