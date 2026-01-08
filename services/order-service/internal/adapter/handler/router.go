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
	userHandler *UserHandler,
	authHandler *AuthHandler,
) *Router {
	// init router
	r := gin.New()

	// group routes
	pb := r.Group("/api/v1")
	// us := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.UserRole, domain.AdminRole))
	ad := pb.Group("/", AuthMiddleware(), RoleMiddleware(domain.AdminRole))

	// public user and auth routes
	pb.POST("/login", authHandler.Login)
	pb.POST("/register", userHandler.RegisterUser)

	// admin user routes
	ad.GET("/users", userHandler.GetUsers)

	return &Router{r}
}

func (r *Router) Serve(conf *config.HTTP) error {
	uri := conf.Host + ":" + conf.Port
	return r.r.Run(uri)
}
