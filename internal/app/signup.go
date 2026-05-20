package app

import (
	"errors"
	authservice "github.com/fidesy-pay/auth-service/internal/pkg/auth-service"
	desc "github.com/fidesy-pay/auth-service/pkg/auth-service"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) SignUp(ctx context.Context, req *desc.SignUpRequest) (*desc.SignUpResponse, error) {
	input, err := authservice.SignUpInputFromRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	token, err := i.authService.SignUp(ctx, input)
	if err != nil {
		if errors.Is(err, authservice.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "authService.Login: %v", err)
	}

	return &desc.SignUpResponse{
		Token: token,
	}, nil
}
