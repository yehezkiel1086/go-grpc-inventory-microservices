package port

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/core/domain"
)

type InventoryRepository interface {
	InitStock(ctx context.Context, inventory *domain.Inventory) (*domain.Inventory, error)
	CheckStock(ctx context.Context, productID int) (*domain.Inventory, error)
	ReduceStock(ctx context.Context, productId int, qty int) (*domain.Inventory, error)
	Restock(ctx context.Context, productId int, qty int) (*domain.Inventory, error)
}

type InventoryService interface {
	InitStock(ctx context.Context, inventory *domain.Inventory) (*domain.Inventory, error)
	CheckStock(ctx context.Context, productID int) (*domain.Inventory, error)
	ReduceStock(ctx context.Context, productId int, qty int) (*domain.Inventory, error)
	Restock(ctx context.Context, productId int, qty int) (*domain.Inventory, error)
}