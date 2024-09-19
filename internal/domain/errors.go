package domain

import "errors"

var (
	ErrStatusStageDone  = errors.New("status stage already done")
	ErrInvalidID        = errors.New("invalid ID")
	ErrInvalidClientID  = errors.New("invalid clientID")
	ErrInvalidCost      = errors.New("invalid cost")
	ErrInvalidWeight    = errors.New("invalid weight")
	ErrStoreTimeExpired = errors.New("store time expired")
	ErrAlreadyPackaged  = errors.New("order already packaged")
	ErrPackageTooHeavy  = errors.New("order too heavy")
)
