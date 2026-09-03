package model

import "errors"

var (
	ErrNotFound      = errors.New("item not found")
	ErrNegativeStock = errors.New("quantity cannot be negative")
)

type Item struct {
	ID        int64  `json:"id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Quantity  int64  `json:"quantity"`
	CreatedAt string `json:"created_at"`
}
