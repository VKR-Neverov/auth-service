package models

import (
	"github.com/dgrijalva/jwt-go"
	desc "github.com/fidesy-pay/auth-service/pkg/auth-service"
)

type TokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	ClientID string `json:"client_id"`
}

func (tc *TokenClaims) Proto() *desc.TokenClaims {
	if tc == nil {
		return nil
	}

	return &desc.TokenClaims{
		Username: tc.Username,
		ClientId: tc.ClientID,
	}
}
