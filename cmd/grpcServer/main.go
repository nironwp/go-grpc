package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/nironwp/grpc/internal/database"
	"github.com/nironwp/grpc/internal/pb"
	"github.com/nironwp/grpc/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

var pl = fmt.Println

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	pl("UnaryInterceptor")
	return handler(ctx, req)
}

func StreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	pl("StreamInterceptor")
	return handler(srv, ss)
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	categoryDb := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDb)

	userDB := database.NewUser(db)
	userService := service.NewUserService(*userDB)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterceptor),
		grpc.StreamInterceptor(StreamInterceptor),
	)
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	pb.RegisterUserServiceServer(grpcServer, userService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		panic(err)
	}

	pl("Server listening")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
