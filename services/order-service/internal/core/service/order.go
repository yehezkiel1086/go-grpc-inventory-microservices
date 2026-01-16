package service

import (
	"context"
	"errors"
	"fmt"

	inventory "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/common/genproto/inventory/protobuf"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/port"
)

type OrderService struct {
	repo port.OrderRepository
	inventoryClient inventory.InventoryServiceClient
	notifRepo port.NotificationRepository
}

func NewOrderService(repo port.OrderRepository, inventoryClient inventory.InventoryServiceClient, notifRepo port.NotificationRepository) *OrderService {
	return &OrderService{
		repo,
		inventoryClient,
		notifRepo,
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
	if err := os.notifRepo.SendNotification(ctx, fmt.Sprintf("%d(%d): new order created", order.ProductID, order.Qty)); err != nil {
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
	order, err := os.repo.UpdatePaymentStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}

	// send notification
	if err := os.notifRepo.SendNotification(ctx, fmt.Sprintf("%d(%d): payment status updated", id, status)); err != nil {
		return nil, err
	}

	return order, nil
}
