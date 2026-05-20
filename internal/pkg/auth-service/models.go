package authservice

import (
	desc "github.com/fidesy-pay/auth-service/pkg/auth-service"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type SignUpInput struct {
	Username string
	Password string
	ClientID uuid.UUID
}

func SignUpInputFromRequest(req *desc.SignUpRequest) (*SignUpInput, error) {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.Username, validation.Required),
		validation.Field(&req.Password, validation.Required),
		validation.Field(&req.ClientId, validation.Required, is.UUIDv4),
	)
	if err != nil {
		return nil, err
	}

	return &SignUpInput{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		ClientID: uuid.MustParse(req.GetClientId()),
	}, nil
}
