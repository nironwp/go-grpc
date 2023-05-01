package service

import (
	"context"

	"github.com/nironwp/grpc/internal/database"
	"github.com/nironwp/grpc/internal/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	UserDB database.User
}

func NewUserService(UserDB database.User) *UserService {
	return &UserService{UserDB: UserDB}
}

func (us *UserService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	user, err := us.UserDB.Create(in.Name, in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	userGrpc := &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return userGrpc, nil
}

func (us *UserService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := us.UserDB.Login(in.Email, in.Password)

	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: token}, err
}
