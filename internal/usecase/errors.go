package usecase

import "errors"

var (
	ErrOrderPickedUp            = errors.New("order picked up")
	ErrOrderDeleted             = errors.New("order deleted")
	ErrOrderStoreTimeNotExpired = errors.New("order store time not expired")

	ErrOrderClientMismatch  = errors.New("order client mismatch")
	ErrOrderIsNotRefundable = errors.New("order is non-refundable")
)
