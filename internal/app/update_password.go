package app

import (
	"context"

	desc "github.com/fidesy-pay/auth-service/pkg/auth-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*desc.UpdatePasswordResponse, error) {
	if err := validateUpdatePasswordRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := i.authService.UpdatePassword(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "authService.UpdatePassword: %v", err)
	}

	return &desc.UpdatePasswordResponse{}, nil
}

func validateUpdatePasswordRequest(req *desc.UpdatePasswordRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Username, validation.Required),
		validation.Field(&req.Password, validation.Required),
	)
}
