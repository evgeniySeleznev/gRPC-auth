package auth

import (
	"context"
	auth_v1 "github.com/evgeniySeleznev/gRPC-auth/pkg/authV1"
	"google.golang.org/grpc"
)

type serverApi struct {
	auth_v1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	auth_v1.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	panic("implement me")
}

func (s *serverApi) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	panic("implement me")
}

func (s *serverApi) IsAdmin(ctx context.Context, req *auth_v1.IsAdminRequest) (*auth_v1.IsAdminResponse, error) {
	panic("implement me")
}
