package clientsconsumer

import (
	"github.com/google/uuid"
	"time"
)

type ClientMessage struct {
	ID                           uuid.UUID  `json:"id,omitempty"`
	Username                     string     `json:"username,omitempty"`
	APIKey                       string     `json:"api_key,omitempty"`
	CreatedAt                    time.Time  `json:"created_at"`
	PhotoURL                     *string    `json:"photo_url,omitempty"`
	Email                        string     `json:"email,omitempty"`
	DeletedAt                    *time.Time `json:"deleted_at,omitempty"`
	IsInvoiceNotificationEnabled bool       `json:"is_invoice_notification_enabled,omitempty"`
}
