package pvz_service

import (
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
	desc "gitlab.ozon.dev/marchenkosasha2/homework/pkg/pvz-service/v1"
)

type Implementation struct {
	usecase usecase.OrderUseCase

	desc.UnimplementedPVZServiceServer
}

func NewImplementation(usecase usecase.OrderUseCase) *Implementation {
	return &Implementation{usecase: usecase}
}
