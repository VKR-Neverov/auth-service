package app

import (
	"context"
	authservice "github.com/fidesy-pay/auth-service/internal/pkg/auth-service"
	"github.com/fidesy-pay/auth-service/internal/pkg/models"
	desc "github.com/fidesy-pay/auth-service/pkg/auth-service"
	"google.golang.org/grpc"
)

type (
	Implementation struct {
		desc.UnimplementedAuthServiceServer

		authService AuthService
	}

	AuthService interface {
		ParseToken(ctx context.Context, token string) (*models.TokenClaims, error)
		Login(ctx context.Context, username, password string) (string, error)
		SignUp(ctx context.Context, input *authservice.SignUpInput) (string, error)
		UpdatePassword(ctx context.Context, username, password string) error
	}
)

func New(
	authService AuthService,
) *Implementation {
	return &Implementation{
		authService: authService,
	}
}

func (i *Implementation) GetDescription() *grpc.ServiceDesc {
	return &desc.AuthService_ServiceDesc
}
