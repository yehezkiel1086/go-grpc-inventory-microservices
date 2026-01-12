package service

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
	inventory "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/common/genproto/inventory/protobuf"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/rabbitmq"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/port"
)

type OrderService struct {
	repo port.OrderRepository
	inventoryClient inventory.InventoryServiceClient
	mq *rabbitmq.Rabbitmq
	q *amqp.Queue
}

func NewOrderService(repo port.OrderRepository, inventoryClient inventory.InventoryServiceClient, mq *rabbitmq.Rabbitmq, q *amqp.Queue) *OrderService {
	return &OrderService{
		repo,
		inventoryClient,
		mq,
		q,
	}
}

func (os *OrderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// reduce stock
	res, err := os.inventoryClient.ReduceStock(ctx, &inventory.ReduceStockReq{
		ProductId: int64(order.ProductID),
		Quantity: int32(order.Qty),
	})
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, errors.New("insufficient stock")
	}

	// calculate total price
	totalPrice := float64(order.Qty) * order.Product.Price
	order.TotalPrice = totalPrice

	// send notification
	if err := os.mq.Publish(ctx, os.q, "new order created"); err != nil {
		return nil, err
	}

	return os.repo.CreateOrder(ctx, order)
}

func (os *OrderService) GetUserOrders(ctx context.Context, userId uint) ([]domain.Order, error) {
	return os.repo.GetUserOrders(ctx, userId)
}

func (os *OrderService) GetOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	return os.repo.GetOrderByID(ctx, id)
}

func (os *OrderService) UpdatePaymentStatus(ctx context.Context, id uint, status domain.OrderStatus) (*domain.Order, error) {
	return os.repo.UpdatePaymentStatus(ctx, id, status)
}
