package repository

import (
	"context"
	"errors"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/core/domain"
)

type InventoryRepository struct {
	db *postgres.DB
}

func NewInventoryRepository(db *postgres.DB) *InventoryRepository {
	return &InventoryRepository{
		db: db,
	}
}

func (repo *InventoryRepository) InitStock(ctx context.Context, inventory *domain.Inventory) (*domain.Inventory, error) {
	db := repo.db.GetDB()

	err := db.WithContext(ctx).Create(inventory).Error
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func (repo *InventoryRepository) CheckStock(ctx context.Context, productID int) (*domain.Inventory, error) {
	db := repo.db.GetDB()

	var inventory domain.Inventory
	if err := db.WithContext(ctx).Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (repo *InventoryRepository) ReduceStock(ctx context.Context, productId int, qty int) (*domain.Inventory, error) {
	db := repo.db.GetDB()

	// get inventory
	var inventory domain.Inventory
	if err := db.WithContext(ctx).Where("product_id = ?", productId).First(&inventory).Error; err != nil {
		return nil, err
	}

	// check if qty sufficient
	if qty > inventory.Qty {
		return nil, errors.New("insufficient stock")
	}

	// reduce stock
	inventory.Qty -= qty
	if err := db.WithContext(ctx).Save(&inventory).Error; err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (repo *InventoryRepository) Restock(ctx context.Context, productId int, qty int) (*domain.Inventory, error) {
	db := repo.db.GetDB()

	// get inventory
	var inventory domain.Inventory
	if err := db.WithContext(ctx).Where("product_id = ?", productId).First(&inventory).Error; err != nil {
		return nil, err
	}

	// restock
	inventory.Qty += qty

	if err := db.WithContext(ctx).Save(&inventory).Error; err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (repo *InventoryRepository) DeleteStock(ctx context.Context, productId int) error {
	db := repo.db.GetDB()
	return db.WithContext(ctx).Where("product_id = ?", productId).Delete(&domain.Inventory{}).Error
}
