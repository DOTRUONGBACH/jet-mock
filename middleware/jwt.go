package middleware

import (
	"context"
	"errors"
	"mock-project/ent"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var jwtKey = uuid.New().String()

func CreateToken(user *ent.User) (string, error) {
	// Tạo claims
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token với secret key
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func jwtAuthFunc(ctx context.Context) (context.Context, error) {
	// Lấy metadata từ context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	// Lấy authorization token từ header
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	// Xác thực JWT token
	token, err := jwt.Parse(authHeader[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid authorization token")
		}

		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization token")
	}

	// truyền thông tin xác thực vào context
	newCtx := context.WithValue(ctx, "user_id", token.Claims.(jwt.MapClaims)["id"].(string))
	newCtx = context.WithValue(newCtx, "email", token.Claims.(jwt.MapClaims)["email"].(string))

	return newCtx, nil
}

func JwtUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newCtx, err := jwtAuthFunc(ctx)
	if err != nil {
		return nil, err
	}

	return handler(newCtx, req)
}

func JwtStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	newCtx, err := jwtAuthFunc(ss.Context())
	if err != nil {
		return err
	}

	wrapped := grpc_middleware.WrapServerStream(ss)
	wrapped.WrappedContext = newCtx

	return handler(srv, wrapped)
}
