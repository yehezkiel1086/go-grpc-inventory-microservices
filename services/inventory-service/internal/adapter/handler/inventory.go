package handler

import (
	"context"

	inventory "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/common/genproto/inventory/protobuf"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/core/port"
)

type InventoryHandler struct {
	inventory.UnimplementedInventoryServiceServer
	svc port.InventoryService
}

func NewInventoryHandler(svc port.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		svc: svc,
	}
}

func (h *InventoryHandler) InitStock(ctx context.Context, req *inventory.InitStockReq) (*inventory.InitStockRes, error) {
	if _, err := h.svc.InitStock(ctx, &domain.Inventory{
		ProductID: uint(req.ProductId),
		Qty:       int(req.Quantity),
	}); err != nil {
		return &inventory.InitStockRes{
			Success: false,
		}, err
	}

	return &inventory.InitStockRes{
		Success: true,
	}, nil
}

func (h *InventoryHandler) CheckStock(ctx context.Context, req *inventory.CheckStockReq) (*inventory.CheckStockRes, error) {
	inv, err := h.svc.CheckStock(ctx, int(req.ProductId))
	if err != nil {
		return nil, err
	}

	return &inventory.CheckStockRes{
		Quantity: int32(inv.Qty),
	}, nil
}

func (h *InventoryHandler) ReduceStock(ctx context.Context, req *inventory.ReduceStockReq) (*inventory.ReduceStockRes, error) {
	if _, err := h.svc.ReduceStock(ctx, int(req.ProductId), int(req.Quantity)); err != nil {
		return &inventory.ReduceStockRes{
			Success: false,
		}, err
	}

	return &inventory.ReduceStockRes{
		Success: true,
	}, nil
}

func (h *InventoryHandler) Restock(ctx context.Context, req *inventory.RestockReq) (*inventory.RestockRes, error) {
	if _, err := h.svc.Restock(ctx, int(req.ProductId), int(req.Quantity)); err != nil {
		return &inventory.RestockRes{
			Success: false,
		}, err
	}

	return &inventory.RestockRes{
		Success: true,
	}, nil
}

func (h *InventoryHandler) DeleteStock(ctx context.Context, req *inventory.DeleteStockReq) (*inventory.DeleteStockRes, error) {
	if err := h.svc.DeleteStock(ctx, int(req.ProductId)); err != nil {
		return &inventory.DeleteStockRes{
			Success: false,
		}, err
	}

	return &inventory.DeleteStockRes{
		Success: true,
	}, nil
}
