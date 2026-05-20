package app

import (
	"context"
	desc "github.com/fidesy-pay/auth-service/pkg/auth-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ValidateToken(ctx context.Context, req *desc.ValidateTokenRequest) (*desc.TokenClaims, error) {
	err := validateValidateTokenRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	token, err := i.authService.ParseToken(ctx, req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "authService.ParseToken: %v", err)
	}

	return token.Proto(), nil
}

func validateValidateTokenRequest(req *desc.ValidateTokenRequest) error {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.Token, validation.Required),
	)
	return err
}
