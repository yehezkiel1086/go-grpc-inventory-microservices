package service

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	inventory "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/common/genproto/inventory/protobuf"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/rabbitmq"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/port"
)

type ProductService struct {
	inventoryClient inventory.InventoryServiceClient
	repo port.ProductRepository
	mq *rabbitmq.Rabbitmq
	q *amqp.Queue
}

func NewProductService(repo port.ProductRepository, inventoryClient inventory.InventoryServiceClient, mq *rabbitmq.Rabbitmq, q *amqp.Queue) *ProductService {
	return &ProductService{
		inventoryClient,
		repo,
		mq,
		q,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, req *domain.CreateProductReq) (*domain.Product, error) {
	// create product
	prod, err := ps.repo.CreateProduct(ctx, &domain.Product{
		Name:  req.Name,
		Price: req.Price,
	})
	if err != nil {
		return nil, err
	}

	// init stock
	if _, err := ps.inventoryClient.InitStock(ctx, &inventory.InitStockReq{
		ProductId: int64(prod.ID),
		Quantity: int32(req.Qty),
	}); err != nil {
		return nil, err
	}

	// send notification
	if err := ps.mq.Publish(ctx, ps.q, "product created"); err != nil {
		return nil, err
	}

	return prod, nil
}

func (ps *ProductService) GetProducts(ctx context.Context) ([]domain.Product, error) {
	return ps.repo.GetProducts(ctx)
}

func (ps *ProductService) GetProductByID(ctx context.Context, id uint) (*domain.ProductRes, error) {
	// get product
	prod, err := ps.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// get quantity (check stock)
	stock, err := ps.inventoryClient.CheckStock(ctx, &inventory.CheckStockReq{
		ProductId: int64(prod.ID),
	})
	if err != nil {
		return nil, err
	}

	return &domain.ProductRes{
		ID:    prod.ID,
		Name:  prod.Name,
		Price: prod.Price,
		Qty:   int(stock.Quantity),
		CreatedAt: prod.CreatedAt,
		UpdatedAt: prod.UpdatedAt,
	}, nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return ps.repo.DeleteProduct(ctx, id)
}
