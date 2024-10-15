package pvz_service

import (
	"context"

	desc "gitlab.ozon.dev/marchenkosasha2/homework/pkg/pvz-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Implementation) OrderList(ctx context.Context, req *desc.OrderListRequest) (*desc.OrderListResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orders, err := s.usecase.OrderList(ctx, int(req.ClientId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	respOrderList := make([]*desc.Order, 0, len(orders.Orders))

	for _, order := range orders.Orders {
		respOrderList = append(respOrderList, &desc.Order{
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

	return &desc.OrderListResponse{Orders: respOrderList}, nil
}
