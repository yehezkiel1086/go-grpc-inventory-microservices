package service

import (
	"context"
	"fmt"

	inventory "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/common/genproto/inventory/protobuf"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/port"
)

type ProductService struct {
	inventoryClient inventory.InventoryServiceClient
	repo port.ProductRepository
	notifRepo port.NotificationRepository
}

func NewProductService(repo port.ProductRepository, inventoryClient inventory.InventoryServiceClient, notifRepo port.NotificationRepository) *ProductService {
	return &ProductService{
		inventoryClient,
		repo,
		notifRepo,
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

	// init stock (grpc)
	if _, err := ps.inventoryClient.InitStock(ctx, &inventory.InitStockReq{
		ProductId: int64(prod.ID),
		Quantity: int32(req.Qty),
	}); err != nil {
		return nil, err
	}

	// send notification (rabbitmq)
	if err := ps.notifRepo.SendNotification(ctx, fmt.Sprintf("%s: product created successfully", prod.Name)); err != nil {
		return nil, err
	}

	return prod, nil
}

func (ps *ProductService) GetProducts(ctx context.Context) ([]domain.ProductRes, error) {
	products, err := ps.repo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	// get quantity (grpc)
	var res []domain.ProductRes
	for _, prod := range products {
		stock, err := ps.inventoryClient.CheckStock(ctx, &inventory.CheckStockReq{
			ProductId: int64(prod.ID),
		})
		if err != nil {
			return nil, err
		}

		res = append(res, domain.ProductRes{
			ID:    prod.ID,
			Name:  prod.Name,
			Price: prod.Price,
			Qty:   int(stock.Quantity),
			CreatedAt: prod.CreatedAt,
			UpdatedAt: prod.UpdatedAt,
		})
	}

	return res, nil
}

func (ps *ProductService) GetProductByID(ctx context.Context, id uint) (*domain.ProductRes, error) {
	// get product
	prod, err := ps.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// get stock quantity (grpc)
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
	// delete product
	if err := ps.repo.DeleteProduct(ctx, id); err != nil {
		return err
	}

	// delete inventory (grpc)
	if _, err := ps.inventoryClient.DeleteStock(ctx, &inventory.DeleteStockReq{
		ProductId: int64(id),
	}); err != nil {
		return err
	}
	
	// send notification (rabbitmq)
	return ps.notifRepo.SendNotification(ctx, fmt.Sprintf("%v: product is deleted", id))
}
