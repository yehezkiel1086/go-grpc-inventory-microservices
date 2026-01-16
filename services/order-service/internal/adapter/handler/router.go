package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(
	conf *config.HTTP,
	productHandler *ProductHandler,
	orderHandler *OrderHandler,
) *Router {
	// init router
	r := gin.New()

	// group routes
	pb := r.Group("/api/v1")
	us := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.UserRole, domain.AdminRole))
	ad := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.AdminRole))

	// public product routes
	pb.GET("/products", productHandler.GetProducts)
	pb.GET("/products/:id", productHandler.GetProductByID)

	// admin product routes
	ad.POST("/products", productHandler.CreateProduct)
	ad.DELETE("/products/:id", productHandler.DeleteProduct)

	// user order routes
	us.POST("/orders", orderHandler.CreateOrder)
	us.GET("/orders", orderHandler.GetUserOrders)
	us.GET("/orders/:id", orderHandler.GetOrderByID)
	us.PUT("/orders/payment/:id", orderHandler.UpdatePaymentStatus)

	return &Router{r}
}

func (r *Router) Serve(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
