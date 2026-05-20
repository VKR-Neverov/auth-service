package clientsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	Consumer struct {
		storage Storage
	}

	Storage interface {
		DeleteUser(ctx context.Context, username string) error
	}
)

func New(storage Storage) *Consumer {
	return &Consumer{
		storage: storage,
	}
}

func (c *Consumer) Consume(ctx context.Context, msg []byte) error {
	var client ClientMessage
	err := json.Unmarshal(msg, &client)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if client.DeletedAt == nil {
		return nil
	}

	err = c.storage.DeleteUser(ctx, client.Username)
	if err != nil {
		return fmt.Errorf("storage.DeleteUser: %w", err)
	}

	return nil
}
