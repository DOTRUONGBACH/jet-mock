package main

import (
	"log"
	"mock-project/database"
	"mock-project/grpc/user-service/handlers"
	"mock-project/grpc/user-service/repositories"

	"mock-project/pb"

	"net"

	"google.golang.org/grpc"
)

func main() {
	client, err := database.NewConnect()
	log.Println(database.NewConnect())
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
	// grpc.UnaryInterceptor(middleware.JwtUnaryServerInterceptor),
	// grpc.StreamInterceptor(middleware.JwtStreamServerInterceptor),
	)

	UserRepository := repositories.UserRepository{
		Client: client,
	}
	userHandler, err := handlers.NewUserHandler(UserRepository)
	if err != nil {
		log.Fatalf("Failed to create userHandler: %v", err)
	}
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	// Lắng nghe các kết nối đến server gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
