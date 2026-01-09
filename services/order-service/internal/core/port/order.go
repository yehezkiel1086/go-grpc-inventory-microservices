package port

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetUserOrders(ctx context.Context, userId uint) ([]domain.Order, error)
	GetOrderByID(ctx context.Context, id uint) (*domain.Order, error)
	UpdatePaymentStatus(ctx context.Context, id uint, status domain.OrderStatus) (*domain.Order, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetUserOrders(ctx context.Context, userId uint) ([]domain.Order, error)
	GetOrderByID(ctx context.Context, id uint) (*domain.Order, error)
	UpdatePaymentStatus(ctx context.Context, id uint, status domain.OrderStatus) (*domain.Order, error)
}
