package dto

import "time"

type OrderDto struct {
	ID         int       `json:"id"`
	ClientID   int       `json:"clientID"`
	StoreUntil time.Time `json:"storeUntil"`
	Status     string    `json:"status"`
	PickUpTime time.Time `json:"pickUpTime,omitempty"`
}

type ListOrdersDTO struct {
	Orders []OrderDto `json:"orders"`
}

type AddOrderRequest struct {
	ID         int       `json:"id"`
	ClientID   int       `json:"clientID"`
	StoreUntil time.Time `json:"storeUntil"`
}
