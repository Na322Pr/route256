package usecase_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository/postgres"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/usecase/mock"
)

func TestOrderUseCase_ReceiveOrderFromCourier(t *testing.T) {
	type args struct {
		req dto.AddOrder
	}

	successStoreTime := time.Now().Add(48 * time.Hour)

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.FacadeMock)
		wantErr  bool
		errValue error
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
			setup: func(facadeMock *mock.FacadeMock) {
				order := dto.OrderDTO{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
					Cost:       1000,
					Weight:     5,
				}

				facadeMock.AddOrderMock.Expect(minimock.AnyContext, order).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ErrorReceiveOrderFromCourier",
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
			setup: func(facadeMock *mock.FacadeMock) {
				order := dto.OrderDTO{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
					Cost:       1000,
					Weight:     5,
				}
				facadeMock.AddOrderMock.Expect(minimock.AnyContext, order).Return(postgres.ErrAlreadyExist)
			},
			wantErr:  true,
			errValue: postgres.ErrAlreadyExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			facadeMock := mock.NewFacadeMock(ctrl)

			tt.setup(facadeMock)

			uc := usecase.NewOrderUseCase(facadeMock)

			err := uc.ReceiveOrderFromCourier(context.Background(), tt.args.req)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
				return
			}
		})
	}
}

func TestOrderUseCase_ReturnOrderToCourier(t *testing.T) {
	type args struct {
		orderID int
	}

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.FacadeMock)
		wantErr  bool
		errValue error
	}{
		{
			name: "Success_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(facadeMock *mock.FacadeMock) {
				successStoreTime := time.Now()

				getOrder := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
				}
				facadeMock.GetOrderByIDMock.Expect(minimock.AnyContext, 10).Return(&getOrder, nil)

				updateOrder := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusDelete],
				}
				facadeMock.UpdateOrderMock.Expect(minimock.AnyContext, updateOrder).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ErrorOrderPickedUp_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(facadeMock *mock.FacadeMock) {
				order := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: time.Now(),
					Status:     domain.OrderStatusMap[domain.OrderStatusPickedUp],
				}

				facadeMock.GetOrderByIDMock.Expect(minimock.AnyContext, 10).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderPickedUp,
		},
		{
			name: "ErrorOrderDeleted_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(facadeMock *mock.FacadeMock) {
				order := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: time.Now(),
					Status:     domain.OrderStatusMap[domain.OrderStatusDelete],
				}

				facadeMock.GetOrderByIDMock.Expect(minimock.AnyContext, 10).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderDeleted,
		},
		{
			name: "ErrorOrderStoreTimeNotExpired_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(facadeMock *mock.FacadeMock) {

				order := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: time.Now().Add(24 * time.Hour),
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
				}

				facadeMock.GetOrderByIDMock.Expect(minimock.AnyContext, 10).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderStoreTimeNotExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			facadeMock := mock.NewFacadeMock(ctrl)
			tt.setup(facadeMock)
			uc := usecase.NewOrderUseCase(facadeMock)

			err := uc.ReturnOrderToCourier(context.Background(), int64(tt.args.orderID))
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestOrderUseCase_GiveOrderToClient(t *testing.T) {
	type args struct {
		orderIDs []int64
	}

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.FacadeMock)
		wantErr  bool
		errValue error
	}{
		{
			name: "Success_GiveOrderToClient",
			args: args{orderIDs: []int64{10}},
			setup: func(facadeMock *mock.FacadeMock) {

				successStoreTime := time.Now().Add(24 * time.Hour)

				orders := &dto.ListOrdersDTO{
					Orders: []dto.OrderDTO{},
				}

				order := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
				}

				orders.Orders = append(orders.Orders, order)
				facadeMock.GetOrdersByIDsMock.Expect(minimock.AnyContext, []int64{10}).Return(orders, nil)
				facadeMock.UpdateOrderMock.Return(nil)

			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := minimock.NewController(t)
			facadeMock := mock.NewFacadeMock(ctrl)
			tt.setup(facadeMock)
			uc := usecase.NewOrderUseCase(facadeMock)

			if err := uc.GiveOrderToClient(context.Background(), tt.args.orderIDs); (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.GiveOrderToClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderUseCase_OrderList(t *testing.T) {
	type args struct {
		clientID int
	}

	successStoreTime := time.Now().Add(48 * time.Hour)

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.FacadeMock)
		want     *dto.ListOrdersDTO
		wantErr  bool
		errValue error
	}{
		{
			name: "SuccessOrderList",
			args: args{
				clientID: 10,
			},
			setup: func(facadeMock *mock.FacadeMock) {

				orders := &dto.ListOrdersDTO{
					Orders: []dto.OrderDTO{},
				}

				for i := 11; i <= 12; i++ {
					order := dto.OrderDTO{
						ID:         int64(i),
						ClientID:   10,
						StoreUntil: successStoreTime,
						Cost:       100 * i,
						Weight:     i,
						Status:     "received",
					}

					orders.Orders = append(orders.Orders, order)
				}

				facadeMock.GetClientOrdersListMock.Expect(minimock.AnyContext, 10).Return(orders, nil)
			},
			want: &dto.ListOrdersDTO{
				Orders: []dto.OrderDTO{
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
			t.Parallel()

			ctrl := minimock.NewController(t)
			facadeMock := mock.NewFacadeMock(ctrl)
			tt.setup(facadeMock)
			uc := usecase.NewOrderUseCase(facadeMock)

			got, err := uc.OrderList(context.Background(), tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.OrderList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrderUseCase_GetRefundFromСlient(t *testing.T) {
	type args struct {
		clientID int
		orderID  int64
	}

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.FacadeMock)
		wantErr  bool
		errValue error
	}{
		{
			name: "Success_GetRefundFromСlient",
			args: args{
				clientID: 10,
				orderID:  11,
			},
			setup: func(facadeMock *mock.FacadeMock) {

				pickUpTime := time.Now()
				order := dto.OrderDTO{
					ID:         11,
					ClientID:   10,
					PickUpTime: sql.NullTime{Time: pickUpTime, Valid: true},
					Status:     domain.OrderStatusMap[domain.OrderStatusPickedUp],
				}

				facadeMock.GetOrderByIDMock.Expect(minimock.AnyContext, 11).Return(&order, nil)
				facadeMock.UpdateOrderMock.Return(nil)

			},
			wantErr: false,
		},
		{
			name: "ErrorOrderClientMismatch_GetRefundFromСlient",
			args: args{
				clientID: 10,
				orderID:  11,
			},
			setup: func(facadeMock *mock.FacadeMock) {

				pickUpTime := time.Now()
				order := dto.OrderDTO{
					ID:         11,
					ClientID:   11,
					PickUpTime: sql.NullTime{Time: pickUpTime, Valid: true},
					Status:     domain.OrderStatusMap[domain.OrderStatusPickedUp],
				}

				facadeMock.GetOrderByIDMock.Expect(minimock.AnyContext, 11).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderClientMismatch,
		},
		{
			name: "ErrorOrderIsNotRefundable_GetRefundFromСlient",
			args: args{
				clientID: 10,
				orderID:  11,
			},
			setup: func(facadeMock *mock.FacadeMock) {
				order := dto.OrderDTO{
					ID:       11,
					ClientID: 10,
					Status:   domain.OrderStatusMap[domain.OrderStatusReceived],
				}

				facadeMock.GetOrderByIDMock.Expect(minimock.AnyContext, 11).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderIsNotRefundable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			facadeMock := mock.NewFacadeMock(ctrl)
			tt.setup(facadeMock)
			uc := usecase.NewOrderUseCase(facadeMock)

			err := uc.GetRefundFromСlient(context.Background(), tt.args.clientID, tt.args.orderID)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestOrderUseCase_RefundList(t *testing.T) {
	type args struct {
		limit  int
		offset int
	}

	successStoreTime := time.Now().Add(48 * time.Hour)

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.FacadeMock)
		want     *dto.ListOrdersDTO
		wantErr  bool
		errValue error
	}{
		{
			name: "SuccessRefundList",
			args: args{
				limit:  0,
				offset: 0,
			},
			setup: func(facadeMock *mock.FacadeMock) {
				refunds := &dto.ListOrdersDTO{
					Orders: []dto.OrderDTO{},
				}

				for i := 11; i <= 12; i++ {
					refund := dto.OrderDTO{
						ID:         int64(i),
						ClientID:   10,
						StoreUntil: successStoreTime,
						Status:     domain.OrderStatusMap[domain.OrderStatusRefunded],
						Cost:       100 * i,
						Weight:     i,
					}

					refunds.Orders = append(refunds.Orders, refund)
				}

				facadeMock.GetRefundsListMock.Expect(minimock.AnyContext, 0, 0).Return(refunds, nil)

			},
			want: &dto.ListOrdersDTO{
				Orders: []dto.OrderDTO{
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
			t.Parallel()

			ctrl := minimock.NewController(t)
			facadeMock := mock.NewFacadeMock(ctrl)
			tt.setup(facadeMock)
			uc := usecase.NewOrderUseCase(facadeMock)

			got, err := uc.RefundList(context.Background(), tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.RefundList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
