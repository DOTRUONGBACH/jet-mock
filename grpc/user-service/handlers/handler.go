package handlers

import (
	"context"
	"log"

	"mock-project/grpc/user-service/models"
	"mock-project/grpc/user-service/repositories"
	"mock-project/middleware"
	"mock-project/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	UserRepository repositories.UserRepository
}

// tao new user
func NewUserHandler(UserRepository repositories.UserRepository) (*UserHandler, error) {
	return &UserHandler{
		UserRepository: UserRepository,
	}, nil
}

func (s *UserHandler) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	// Kiểm tra xem email đã được sử dụng chưa
	user1, err := s.UserRepository.FindByEmail(ctx, req.User.Email)
	log.Println("check", user1)
	if err == nil {
		return &pb.SignupResponse{
			Success: false,
			Message: "Email has been used",
		}, status.Errorf(codes.AlreadyExists, "Email has been used")
	}

	// Tạo mới một user trong database
	newUser := &models.User{
		Email:       req.User.Email,
		FullName:    req.User.FullName,
		PhoneNumber: req.User.PhoneNumber,
		IDCard:      req.User.IdCard,
		DateOfBirth: req.User.DateOfBirth,
		Password:    req.Password,
		MemberCard:  int(req.User.MemberCard),
	}
	_, err2 := s.UserRepository.Create(newUser)
	if err2 != nil {
		return nil, err2
	}

	return &pb.SignupResponse{
		Success: true,
		Message: "successful",
	}, nil
}

func (s *UserHandler) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	user, err := s.UserRepository.FindByEmail(ctx, req.Email)

	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Kiểm tra mật khẩu
	check := s.UserRepository.CheckUserPassWord(ctx, req.Email, req.Password)
	log.Println(check)
	if check == false {
		return nil, status.Error(codes.NotFound, "wrong passs")
	}

	// tạo JWT token và trả về trong phản hồi SigninResponse
	token, err := middleware.CreateToken(user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create token")
	}

	return &pb.SigninResponse{
		Success: true,
		Token:   token,
		Message: "Login successful",
	}, nil
}
