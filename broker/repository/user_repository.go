package repository

import (
	"broker/proto"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserRepository interface {
	Register(*proto.RegisterRequest) error
	Login(*proto.LoginRequest) (*proto.TokenResponse, error)
	RefreshToken(*proto.RefreshTokenRequest) (*proto.TokenResponse, error)
	GetUserByID(ID int) (*proto.User, error)
}

type UserRepositoryImpl struct {
	client proto.AuthServiceClient
}

func NewUserRepository() *UserRepositoryImpl {
	conn, err := grpc.NewClient("user_service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to connect: %v", err)
	}
	logrus.Info("Connected to user service")

	client := proto.NewAuthServiceClient(conn)

	return &UserRepositoryImpl{client: client}
}

func (u *UserRepositoryImpl) Register(payload *proto.RegisterRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := u.client.Register(ctx, &proto.RegisterRequest{
		Payload: &proto.RegisterPayload{
			FullName: payload.Payload.FullName,
			Email:    payload.Payload.Email,
			Password: payload.Payload.Password,
		},
	})
	if err != nil {
		logrus.Error("Failed to register: ", err)
		return err
	}

	logrus.Info("Registered successfully")
	return nil
}

func (u *UserRepositoryImpl) Login(payload *proto.LoginRequest) (*proto.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.Login(ctx, &proto.LoginRequest{
		Payload: &proto.LoginPayload{
			Email:    payload.Payload.Email,
			Password: payload.Payload.Password,
		},
	})
}

func (u *UserRepositoryImpl) RefreshToken(payload *proto.RefreshTokenRequest) (*proto.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.RefreshToken(ctx, &proto.RefreshTokenRequest{
		Payload: &proto.RefreshTokenPayload{
			RefreshToken: payload.Payload.RefreshToken,
		},
	})
}

func (u *UserRepositoryImpl) GetUserByID(ID int) (*proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.GetUserByID(ctx, &proto.GetUserRequest{
		UserId: int64(ID)},
	)
}
