package usecase

import (
	"fmt"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase/mock"
)

func TestOrderUseCase_ReceiveOrderFromCourier(t *testing.T) {
	type args struct {
		req dto.AddOrder
	}

	successStoreTime := time.Now().AddDate(0, 0, 2)

	ctrl := minimock.NewController(t)
	repoMock := mock.NewRepositoryMock(ctrl)

	uc := &OrderUseCase{
		repo: repoMock,
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "SuccessReceiveOrderFromCourier",
			args: args{
				req: dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
					Packages:   []string{"unknown", "unknown"},
				},
			},
			setup: func() {
				order, _ := domain.NewOrder(dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
					Packages:   []string{"unknown", "unknown"},
				}, nil, nil)
				fmt.Println(order)
				repoMock.AddOrderMock.Expect(order).Return(nil)
				repoMock.UpdateMock.Expect().Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.setup()

		t.Run(tt.name, func(t *testing.T) {
			if err := uc.ReceiveOrderFromCourier(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.ReceiveOrderFromCourier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderUseCase_ReturnOrderToCourier(t *testing.T) {
	type args struct {
		orderID int
	}

	ctrl := minimock.NewController(t)
	repoMock := mock.NewRepositoryMock(ctrl)

	uc := &OrderUseCase{
		repo: repoMock,
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "SuccessReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func() {
				var order domain.Order
				order.SetClientID(10)
				order.SetStoreUntil(time.Now())
				order.SetStatus(domain.OrderStatusReceived)

				repoMock.GetOrderByIDMock.Expect(10).Return(&order, nil)
				repoMock.UpdateMock.Expect().Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.setup()

		t.Run(tt.name, func(t *testing.T) {
			if err := uc.ReturnOrderToCourier(tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.ReturnOrderToCourier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderUseCase_GiveOrderToClient(t *testing.T) {
	type args struct {
		orderIDs []int
	}

	ctrl := minimock.NewController(t)
	repoMock := mock.NewRepositoryMock(ctrl)

	uc := &OrderUseCase{
		repo: repoMock,
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "SuccessReturnOrderToCourier",
			args: args{orderIDs: []int{10}},
			setup: func() {
				var orders []*domain.Order

				var order domain.Order
				order.SetClientID(10)
				order.SetStoreUntil(time.Now().AddDate(0, 0, 1))
				order.SetStatus(domain.OrderStatusReceived)

				orders = append(orders, &order)

				repoMock.GetOrdersByIDMock.Expect([]int{10}).Return(orders, nil)
				repoMock.UpdateMock.Expect().Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			if err := uc.GiveOrderToClient(tt.args.orderIDs); (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.GiveOrderToClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderUseCase_OrderList(t *testing.T) {
	type args struct {
		clientID int
	}

	successStoreTime := time.Now().AddDate(0, 0, 2)

	ctrl := minimock.NewController(t)
	repoMock := mock.NewRepositoryMock(ctrl)

	uc := &OrderUseCase{
		repo: repoMock,
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		want    *dto.ListOrdersDTO
		wantErr bool
	}{
		{
			name: "SuccessOrderList",
			args: args{
				clientID: 10,
			},
			setup: func() {
				var orders []*domain.Order

				for i := 11; i <= 12; i++ {
					order, _ := domain.NewOrder(dto.AddOrder{
						ID:         i,
						ClientID:   10,
						StoreUntil: successStoreTime,
						Cost:       100 * i,
						Weight:     i,
						Packages:   []string{},
					})

					orders = append(orders, order)
				}

				repoMock.GetClientOrdersListMock.Expect(10).Return(orders, nil)
			},
			want: &dto.ListOrdersDTO{
				Orders: []dto.OrderDto{
					{
						ID:         11,
						ClientID:   10,
						StoreUntil: successStoreTime,
						Status:     "received",
						Cost:       1100,
						Weight:     11,
					},
					{
						ID:         12,
						ClientID:   10,
						StoreUntil: successStoreTime,
						Status:     "received",
						Cost:       1200,
						Weight:     12,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := uc.OrderList(tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.OrderList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestOrderUseCase_GetRefundFromСlient(t *testing.T) {
	type args struct {
		clientID int
		orderID  int
	}

	ctrl := minimock.NewController(t)
	repoMock := mock.NewRepositoryMock(ctrl)

	uc := &OrderUseCase{
		repo: repoMock,
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "SuccessGetRefundFromСlient",
			args: args{
				clientID: 10,
				orderID:  11,
			},
			setup: func() {
				var order domain.Order
				order.SetID(11)
				order.SetClientID(10)
				order.SetPickUpTime()
				order.SetStatus(domain.OrderStatusPickedUp)

				repoMock.GetOrderByIDMock.Expect(11).Return(&order, nil)
				repoMock.UpdateMock.Expect().Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			if err := uc.GetRefundFromСlient(tt.args.clientID, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.GetRefundFromСlient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderUseCase_RefundList(t *testing.T) {
	type args struct {
		limit  int
		offset int
	}

	successStoreTime := time.Now().AddDate(0, 0, 2)

	ctrl := minimock.NewController(t)
	repoMock := mock.NewRepositoryMock(ctrl)

	uc := &OrderUseCase{
		repo: repoMock,
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		want    *dto.ListOrdersDTO
		wantErr bool
	}{
		{
			name: "SuccessRefundList",
			args: args{
				limit:  0,
				offset: 0,
			},
			setup: func() {
				var refunds []*domain.Order

				for i := 11; i <= 12; i++ {
					var refund domain.Order
					refund.SetID(i)
					refund.SetClientID(10)
					refund.SetStoreUntil(successStoreTime)
					refund.SetStatus(domain.OrderStatusRefunded)
					refund.SetCost(100 * i)
					refund.SetWeight(i)

					refunds = append(refunds, &refund)
				}

				repoMock.GetRefundsListMock.Expect(0, 0).Return(refunds, nil)

			},
			want: &dto.ListOrdersDTO{
				Orders: []dto.OrderDto{
					{
						ID:         11,
						ClientID:   10,
						StoreUntil: successStoreTime,
						Status:     "refunded",
						Cost:       1100,
						Weight:     11,
					},
					{
						ID:         12,
						ClientID:   10,
						StoreUntil: successStoreTime,
						Status:     "refunded",
						Cost:       1200,
						Weight:     12,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := uc.RefundList(tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.RefundList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}
