package domain

import "errors"

var (
	ErrStatusStageDone = errors.New("status stage already done")
)
