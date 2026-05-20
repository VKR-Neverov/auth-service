package app

import (
	"context"
	"errors"
	"github.com/fidesy-pay/auth-service/internal/pkg/storage"
	desc "github.com/fidesy-pay/auth-service/pkg/auth-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	err := validateLoginRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	token, err := i.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "authService.Login: %v", err)
	}

	return &desc.LoginResponse{
		Token: token,
	}, nil
}

func validateLoginRequest(req *desc.LoginRequest) error {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.Username, validation.Required),
		validation.Field(&req.Password, validation.Required),
	)
	return err
}
