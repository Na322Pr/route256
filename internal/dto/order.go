package dto

import (
	"time"
)

type OrderDTO struct {
	ID         int        `json:"id" db:"order_id"`
	ClientID   int        `json:"clientID" db:"client_id"`
	StoreUntil time.Time  `json:"storeUntil" db:"store_until"`
	Status     string     `json:"status" db:"status"`
	Cost       int        `json:"cost" db:"cost"`
	Weight     int        `json:"weight" db:"weight"`
	Packages   []string   `json:"packages" db:"packages"`
	PickUpTime *time.Time `json:"pickUpTime,omitempty" db:"pick_up_time"`
}

type ListOrdersDTO struct {
	Orders []OrderDTO `json:"orders"`
}

type AddOrder struct {
	ID         int       `json:"id"`
	ClientID   int       `json:"clientID"`
	StoreUntil time.Time `json:"storeUntil"`
	Cost       int       `json:"cost"`
	Weight     int       `json:"weight"`
	Packages   []string  `json:"packages"`
}
