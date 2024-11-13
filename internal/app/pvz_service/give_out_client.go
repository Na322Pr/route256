package pvz_service

import (
	"context"

	desc "github.com/Na322Pr/route256/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GiveOutClient(ctx context.Context, req *desc.GiveOutClientRequest) (*desc.GiveOutClientResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.usecase.GiveOrderToClient(ctx, req.OrdersIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.GiveOutClientResponse{}, nil
}
