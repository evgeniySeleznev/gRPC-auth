package auth

import (
	"context"
	"errors"
	"github.com/evgeniySeleznev/gRPC-auth/internal/services/auth"
	auth_v1 "github.com/evgeniySeleznev/gRPC-auth/pkg/authV1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//место использования интерфейса — интерфейс в сигнатуре функции
//место реализации — где описывается соответствие интерфейсу, без использования

//здесь содержится основная бизнес-логика
//не в хендлерах ниже, а здесь в сервисном слое

type Auth interface {
	Get(ctx context.Context, email string, password string, appID int32) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverApi struct {
	auth_v1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	auth_v1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

const (
	emptyValue = 0
)

func (s *serverApi) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Get(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth_v1.GetResponse{Token: token}, nil
}

func (s *serverApi) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	if err := validateLRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth_v1.CreateResponse{
		Id: userID,
	}, nil
}

func (s *serverApi) IsAdmin(ctx context.Context, req *auth_v1.IsAdminRequest) (*auth_v1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth_v1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

// валидация логина
func validateLogin(req *auth_v1.GetRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

// валидация регистрации
func validateLRegister(req *auth_v1.CreateRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// валидация проверки на админа
func validateIsAdmin(req *auth_v1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "empty user id")
	}

	return nil
}
