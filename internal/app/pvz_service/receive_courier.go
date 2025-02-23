package pvz_service

import (
	"context"

	"github.com/Na322Pr/route256/internal/dto"
	desc "github.com/Na322Pr/route256/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) ReceiveCourier(ctx context.Context, req *desc.ReceiveCourierRequest) (*desc.ReceiveCourierResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	addOrderDTO := dto.AddOrder{
		ID:         req.OrderId,
		ClientID:   int(req.ClientId),
		StoreUntil: req.StoreUntil.AsTime(),
		Cost:       int(req.Cost),
		Weight:     int(req.Weight),
		Packages:   req.Packages,
	}

	err := s.usecase.ReceiveOrderFromCourier(ctx, addOrderDTO)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.ReceiveCourierResponse{}, nil
}
