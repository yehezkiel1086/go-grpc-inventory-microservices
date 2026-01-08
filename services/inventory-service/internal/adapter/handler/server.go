package handler

import (
	"fmt"
	"net"

	inventory "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/common/genproto/inventory/protobuf"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/adapter/config"
	"google.golang.org/grpc"
)

type Server struct {
	s *grpc.Server
}

func NewServer(
	invHandler *InventoryHandler,
) *Server {
	s := grpc.NewServer()

	inventory.RegisterInventoryServiceServer(
		s,
		invHandler,
	)

	return &Server{
		s: s,
	}
}

func (s *Server) Run(conf *config.GRPC) error {
	uri := fmt.Sprintf("%s:%s", conf.Host, conf.Port)

	lis, err := net.Listen("tcp", uri)
	if err != nil {
		return err
	}

	return s.s.Serve(lis)
}
