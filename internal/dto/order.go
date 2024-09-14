package dto

import (
	"time"
)

type OrderDto struct {
	ID         int       `json:"id"`
	ClientID   int       `json:"clientID"`
	StoreUntil time.Time `json:"storeUntil"`
	Status     string    `json:"status"`
	Cost       int       `json:"cost"`
	Weight     int       `json:"weight"`
	Packages   []string  `json:"packages"`
	PickUpTime time.Time `json:"pickUpTime,omitempty"`
}

type ListOrdersDTO struct {
	Orders []OrderDto `json:"orders"`
}

type AddOrder struct {
	ID         int       `json:"id"`
	ClientID   int       `json:"clientID"`
	StoreUntil time.Time `json:"storeUntil"`
	Cost       int       `json:"cost"`
	Weight     int       `json:"weight"`
	Packages   []string  `json:"packages"`
}
