package authservice

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fidesy-pay/auth-service/internal/config"
	"github.com/fidesy-pay/auth-service/internal/pkg/models"
	"github.com/fidesy-pay/auth-service/internal/pkg/storage"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type (
	Service struct {
		storage Storage
	}

	Storage interface {
		CreateUser(ctx context.Context, user *models.User) error
		GetUser(ctx context.Context, username string) (*models.User, error)
		UpdateUser(ctx context.Context, user *models.User) error
	}
)

func New(
	storage Storage,
) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.storage.GetUser(ctx, username)
	if err != nil {
		return "", fmt.Errorf("storage.GetUser: %w", err)
	}

	if !checkPasswordHash(password, user.Password) {
		return "", storage.ErrNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.Get(config.TokenTTL).(time.Duration)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: user.Username,
		ClientID: user.ClientID.String(),
	})

	return token.SignedString([]byte(
		config.Get(config.SigningKey).(string),
	))
}

func (s *Service) SignUp(ctx context.Context, input *SignUpInput) (string, error) {
	_, err := s.storage.GetUser(ctx, input.Username)
	if err == nil {
		return "", ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return "", fmt.Errorf("hashPassword: %w", err)
	}

	user := &models.User{
		Username: input.Username,
		Password: hashedPassword,
		ClientID: input.ClientID,
	}

	err = s.storage.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("storage.CreateUser: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.Get(config.TokenTTL).(time.Duration)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: user.Username,
		ClientID: user.ClientID.String(),
	})

	return token.SignedString([]byte(
		config.Get(config.SigningKey).(string),
	))
}

func (s *Service) UpdatePassword(ctx context.Context, username, password string) error {
	user, err := s.storage.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("storage.GetUser: %w", err)
	}

	user.Password, err = hashPassword(password)
	if err != nil {
		return fmt.Errorf("hashPassword: %w", err)
	}

	err = s.storage.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("storage.UpdateUser: %w", err)
	}

	return nil
}

func (s *Service) ParseToken(_ context.Context, accessToken string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(config.Get(config.SigningKey).(string)), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims, nil
}

func hashPassword(password string) (string, error) {
	// Create a new SHA-256 hash
	hasher := sha256.New()

	// Write the password bytes to the hash
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}

	// Get the final hash as a byte slice
	hashBytes := hasher.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

func checkPasswordHash(password, hash string) bool {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return false
	}

	return hashedPassword == hash
}
