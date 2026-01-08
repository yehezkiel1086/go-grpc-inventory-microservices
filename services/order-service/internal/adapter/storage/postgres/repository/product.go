package repository

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
)

type ProductRepository struct {
	db *postgres.DB
}

func NewProductRepository(db *postgres.DB) *ProductRepository {
	return &ProductRepository{
		db,
	}
}

func (pr *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	db := pr.db.GetDB()
	if err := db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *ProductRepository) GetProducts(ctx context.Context) ([]domain.Product, error) {
	db := pr.db.GetDB()

	var products []domain.Product
	if err := db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *ProductRepository) GetProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	db := pr.db.GetDB()

	var product *domain.Product
	if err := db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	db := pr.db.GetDB()

	if err := db.WithContext(ctx).Delete(&domain.Product{}, id).Error; err != nil {
		return err
	}

	return nil
}
