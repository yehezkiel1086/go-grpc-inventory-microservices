package service

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/core/port"
)

type InventoryService struct {
	repo port.InventoryRepository
}

func NewInventoryService(repo port.InventoryRepository) *InventoryService {
	return &InventoryService{
		repo,
	}
}

func (s *InventoryService) InitStock(ctx context.Context, inventory *domain.Inventory) (*domain.Inventory, error) {
	return s.repo.InitStock(ctx, inventory)
}

func (s *InventoryService) CheckStock(ctx context.Context, productID int) (*domain.Inventory, error) {
	return s.repo.CheckStock(ctx, productID)
}

func (s *InventoryService) ReduceStock(ctx context.Context, productId int, qty int) (*domain.Inventory, error) {
	return s.repo.ReduceStock(ctx, productId, qty)
}

func (s *InventoryService) Restock(ctx context.Context, productId int, qty int) (*domain.Inventory, error) {
	return s.repo.Restock(ctx, productId, qty)
}
