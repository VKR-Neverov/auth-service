package models

import "github.com/google/uuid"

type User struct {
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
	ClientID uuid.UUID `json:"client_id,omitempty"`
}
