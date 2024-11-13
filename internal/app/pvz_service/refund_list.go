package pvz_service

import (
	"context"

	desc "github.com/Na322Pr/route256/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Implementation) RefundList(ctx context.Context, req *desc.RefundListRequest) (*desc.RefundListResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orders, err := s.usecase.RefundList(ctx, int(*req.Limit), int(*req.Offset))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	respRefundList := make([]*desc.Order, 0, len(orders.Orders))

	for _, order := range orders.Orders {
		respRefundList = append(respRefundList, &desc.Order{
			Id:         order.ID,
			ClientId:   int32(order.ClientID),
			StoreUntil: timestamppb.New(order.StoreUntil),
			Status:     order.Status,
			Cost:       int32(order.Cost),
			Weight:     int32(order.Weight),
			Packages:   order.Packages,
			PickUpTime: timestamppb.New(order.StoreUntil),
		})
	}

	return &desc.RefundListResponse{Orders: respRefundList}, nil
}
