package pvz_service

import (
	"github.com/Na322Pr/route256/internal/usecase"
	desc "github.com/Na322Pr/route256/pkg/pvz-service/v1"
)

type Implementation struct {
	usecase usecase.OrderUseCase

	desc.UnimplementedPVZServiceServer
}

func NewImplementation(usecase usecase.OrderUseCase) *Implementation {
	return &Implementation{usecase: usecase}
}
