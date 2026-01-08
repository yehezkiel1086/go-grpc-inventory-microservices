package service

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/port"
)

type ProductService struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) *ProductService {
	return &ProductService{
		repo,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	return ps.repo.CreateProduct(ctx, product)
}

func (ps *ProductService) GetProducts(ctx context.Context) ([]domain.Product, error) {
	return ps.repo.GetProducts(ctx)
}

func (ps *ProductService) GetProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	return ps.repo.GetProductByID(ctx, id)
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return ps.repo.DeleteProduct(ctx, id)
}