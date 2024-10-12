package pvz_service

import (
	"context"

	desc "gitlab.ozon.dev/marchenkosasha2/homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) ReturnCourier(ctx context.Context, req *desc.ReturnCourierRequest) (*desc.ReturnCourierResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.usecase.ReturnOrderToCourier(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.ReturnCourierResponse{}, nil
}
