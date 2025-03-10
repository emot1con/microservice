package grpc

import (
	"context"
	"net"
	"order/cmd/db"
	"order/proto"
	"order/repository"
	"order/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type OrderGRPCServer struct {
	service *service.OrderService
	proto.UnimplementedOrderServiceServer
}

func NewOrderGRPCServer(service *service.OrderService) *OrderGRPCServer {
	return &OrderGRPCServer{
		service: service,
	}
}

func (u *OrderGRPCServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	logrus.Info("create order")
	order, err := u.service.CreateOrder(req)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *OrderGRPCServer) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.OrderResponse, error) {
	order, err := u.service.GetOrderByID(req)
	if err != nil {
		return nil, err
	}

	return &proto.OrderResponse{
		Order: order,
	}, nil
}

func GRPCListen() {
	DB, err := db.Connect()
	if err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}

	orderRepo := repository.NewOrderRepositoryImpl()
	orderItemRepo := repository.NewOrderItemsRepositoryImpl()

	orderService := service.NewOrderItemService(DB, orderRepo, orderItemRepo)
	orderGRPC := NewOrderGRPCServer(orderService)

	conn, err := net.Listen("tcp", ":30001")
	if err != nil {
		logrus.Fatalf("failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterOrderServiceServer(srv, orderGRPC)
	logrus.Infof("gRPC Server started on port 30001")

	if err := srv.Serve(conn); err != nil {
		logrus.Fatalf("error when connect to gRPC Server  %v", err)
	}
}
