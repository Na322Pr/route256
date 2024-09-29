package postgres

import "errors"

var (
	ErrAlreadyExist  = errors.New("order already exist")
	ErrOrderNotFound = errors.New("order not found")
)
