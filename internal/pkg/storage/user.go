package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/fidesy-pay/auth-service/internal/pkg/models"
)

var (
	ErrNotFound = errors.New("entity not found")
)

func (s *Storage) GetUser(ctx context.Context, username string) (*models.User, error) {
	user := new(models.User)

	ok, err := s.redis.Get(ctx, username, &user)
	if err != nil {
		return nil, fmt.Errorf("redis.Get: %w", err)
	}

	if !ok {
		return nil, ErrNotFound
	}

	return user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	err := s.redis.Set(ctx, user.Username, user, 0)
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *models.User) error {
	err := s.redis.Set(ctx, user.Username, user, 0)
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, username string) error {
	err := s.redis.Delete(ctx, username)
	if err != nil {
		return fmt.Errorf("redis.Delete: %w", err)
	}

	return nil
}
