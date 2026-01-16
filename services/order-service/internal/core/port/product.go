package port

import (
	"context"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProducts(ctx context.Context) ([]domain.Product, error)
	GetProductByID(ctx context.Context, id uint) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *domain.CreateProductReq) (*domain.Product, error)
	GetProducts(ctx context.Context) ([]domain.ProductRes, error)
	GetProductByID(ctx context.Context, id uint) (*domain.ProductRes, error)
	DeleteProduct(ctx context.Context, id uint) error
}
