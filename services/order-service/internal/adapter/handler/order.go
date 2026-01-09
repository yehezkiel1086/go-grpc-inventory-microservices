package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/port"
)

type OrderHandler struct {
	svc port.OrderService
}

func NewOrderHandler(svc port.OrderService) *OrderHandler {
	return &OrderHandler{
		svc,
	}
}

type CreateOrderReq struct {
	ProductID uint `json:"product_id" binding:"required"`
	Qty       uint `json:"qty" binding:"required"`
}

func (oh *OrderHandler) CreateOrder(c *gin.Context) {
	var req *CreateOrderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := user.(*domain.JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user claims"})
		return
	}

	order, err := oh.svc.CreateOrder(c.Request.Context(), &domain.Order{
		UserID:    claims.ID,
		ProductID: req.ProductID,
		Qty:       req.Qty,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (oh *OrderHandler) GetUserOrders(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := user.(*domain.JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user claims"})
		return
	}

	orders, err := oh.svc.GetUserOrders(c.Request.Context(), claims.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (oh *OrderHandler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := oh.svc.GetOrderByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (oh *OrderHandler) UpdatePaymentStatus(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req struct {
		Status domain.OrderStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := oh.svc.UpdatePaymentStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
