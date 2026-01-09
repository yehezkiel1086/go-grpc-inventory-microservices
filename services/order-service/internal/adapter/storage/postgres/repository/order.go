package repository

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
)

type OrderRepository struct {
	db *postgres.DB
}

func NewOrderRepository(db *postgres.DB) *OrderRepository {
	return &OrderRepository{
		db,
	}
}

func (or *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	db := or.db.GetDB()
	if err := db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (or *OrderRepository) GetUserOrders(ctx context.Context, userId uint) ([]domain.Order, error) {
	db := or.db.GetDB()

	var orders []domain.Order
	if err := db.WithContext(ctx).Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *OrderRepository) GetOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	db := or.db.GetDB()

	var order *domain.Order
	if err := db.WithContext(ctx).Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (or *OrderRepository) UpdatePaymentStatus(ctx context.Context, id uint, status domain.OrderStatus) (*domain.Order, error) {
	db := or.db.GetDB()

	var order *domain.Order
	if err := db.WithContext(ctx).Model(&order).Where("id = ?", id).Update("status", status).Error; err != nil {
		return nil, err
	}

	return order, nil
}
