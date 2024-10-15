package pvz_service

import (
	"context"

	desc "gitlab.ozon.dev/marchenkosasha2/homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) RefundClient(ctx context.Context, req *desc.RefundClientRequest) (*desc.RefundClientResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.usecase.GetRefundFromСlient(ctx, int(req.ClientId), req.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.RefundClientResponse{}, nil
}
