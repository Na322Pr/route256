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
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/kafka/event"
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
		name  string
		args  args
		setup func(
			*mock.OrderRepoFacadeMock,
			*mock.EventLogProducerFacadeMock,
			*mock.OrderCacheFacadeMock,
		)
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
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				prodMock *mock.EventLogProducerFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {
				order := dto.OrderDTO{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
					Cost:       1000,
					Weight:     5,
					PickUpTime: sql.NullTime{Valid: true},
				}

				repoMock.AddOrderMock.Expect(minimock.AnyContext, order).Return(nil)

				prodMock.ProduceEventMock.Expect(order, event.EventTypeReceive).Return(nil)
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
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				prodMock *mock.EventLogProducerFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {
				order := dto.OrderDTO{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
					Cost:       1000,
					Weight:     5,
					PickUpTime: sql.NullTime{Valid: true},
				}
				repoMock.AddOrderMock.Expect(minimock.AnyContext, order).Return(postgres.ErrAlreadyExist)
			},
			wantErr:  true,
			errValue: postgres.ErrAlreadyExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			repoMock := mock.NewOrderRepoFacadeMock(ctrl)
			prodMock := mock.NewEventLogProducerFacadeMock(ctrl)
			cacheMock := mock.NewOrderCacheFacadeMock(ctrl)

			tt.setup(repoMock, prodMock, cacheMock)

			uc := usecase.NewOrderUseCase(repoMock, prodMock, cacheMock)

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
		name  string
		args  args
		setup func(
			*mock.OrderRepoFacadeMock,
			*mock.OrderCacheFacadeMock,
		)
		wantErr  bool
		errValue error
	}{
		{
			name: "Success_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {
				successStoreTime := time.Now()

				getOrder := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
				}
				repoMock.GetOrderByIDMock.Expect(minimock.AnyContext, 10).Return(&getOrder, nil)

				updateOrder := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: successStoreTime,
					Status:     domain.OrderStatusMap[domain.OrderStatusDelete],
					PickUpTime: sql.NullTime{Valid: true},
				}
				repoMock.UpdateOrderMock.Expect(minimock.AnyContext, updateOrder).Return(nil)

				cacheMock.GetMock.Expect(10).Return(&dto.OrderDTO{}, false)
				cacheMock.SetMock.Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ErrorOrderPickedUp_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {
				order := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: time.Now(),
					Status:     domain.OrderStatusMap[domain.OrderStatusPickedUp],
				}

				repoMock.GetOrderByIDMock.Expect(minimock.AnyContext, 10).Return(&order, nil)

				cacheMock.GetMock.Expect(10).Return(&dto.OrderDTO{}, false)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderPickedUp,
		},
		{
			name: "ErrorOrderDeleted_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {
				order := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: time.Now(),
					Status:     domain.OrderStatusMap[domain.OrderStatusDelete],
				}

				repoMock.GetOrderByIDMock.Expect(minimock.AnyContext, 10).Return(&order, nil)

				cacheMock.GetMock.Expect(10).Return(&dto.OrderDTO{}, false)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderDeleted,
		},
		{
			name: "ErrorOrderStoreTimeNotExpired_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {
				order := dto.OrderDTO{
					ClientID:   10,
					StoreUntil: time.Now().Add(24 * time.Hour),
					Status:     domain.OrderStatusMap[domain.OrderStatusReceived],
				}

				repoMock.GetOrderByIDMock.Expect(minimock.AnyContext, 10).Return(&order, nil)
				cacheMock.GetMock.Expect(10).Return(&dto.OrderDTO{}, false)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderStoreTimeNotExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			repoMock := mock.NewOrderRepoFacadeMock(ctrl)
			prodMock := mock.NewEventLogProducerFacadeMock(ctrl)
			cacheMock := mock.NewOrderCacheFacadeMock(ctrl)

			tt.setup(repoMock, cacheMock)
			uc := usecase.NewOrderUseCase(repoMock, prodMock, cacheMock)

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
		name  string
		args  args
		setup func(
			*mock.OrderRepoFacadeMock,
			*mock.EventLogProducerFacadeMock,
			*mock.OrderCacheFacadeMock,
		)
		wantErr  bool
		errValue error
	}{
		{
			name: "Success_GiveOrderToClient",
			args: args{orderIDs: []int64{10}},
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				prodMock *mock.EventLogProducerFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {

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
				repoMock.GetOrdersByIDsMock.Expect(minimock.AnyContext, []int64{10}).Return(orders, nil)
				repoMock.UpdateOrderMock.Return(nil)

				prodMock.ProduceEventMock.Return(nil)

				cacheMock.SetMock.Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := minimock.NewController(t)
			repoMock := mock.NewOrderRepoFacadeMock(ctrl)
			prodMock := mock.NewEventLogProducerFacadeMock(ctrl)
			cacheMock := mock.NewOrderCacheFacadeMock(ctrl)

			tt.setup(repoMock, prodMock, cacheMock)
			uc := usecase.NewOrderUseCase(repoMock, prodMock, cacheMock)

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
		setup    func(*mock.OrderRepoFacadeMock)
		want     *dto.ListOrdersDTO
		wantErr  bool
		errValue error
	}{
		{
			name: "SuccessOrderList",
			args: args{
				clientID: 10,
			},
			setup: func(repoMock *mock.OrderRepoFacadeMock) {

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

				repoMock.GetClientOrdersListMock.Expect(minimock.AnyContext, 10).Return(orders, nil)
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
			repoMock := mock.NewOrderRepoFacadeMock(ctrl)
			prodMock := mock.NewEventLogProducerFacadeMock(ctrl)
			cacheMock := mock.NewOrderCacheFacadeMock(ctrl)

			tt.setup(repoMock)
			uc := usecase.NewOrderUseCase(repoMock, prodMock, cacheMock)

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
		name  string
		args  args
		setup func(
			*mock.OrderRepoFacadeMock,
			*mock.EventLogProducerFacadeMock,
			*mock.OrderCacheFacadeMock,
		)
		wantErr  bool
		errValue error
	}{
		{
			name: "Success_GetRefundFromСlient",
			args: args{
				clientID: 10,
				orderID:  11,
			},
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				prodMock *mock.EventLogProducerFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {
				pickUpTime := time.Now()
				order := dto.OrderDTO{
					ID:         11,
					ClientID:   10,
					PickUpTime: sql.NullTime{Time: pickUpTime, Valid: true},
					Status:     domain.OrderStatusMap[domain.OrderStatusPickedUp],
				}

				repoMock.GetOrderByIDMock.Expect(minimock.AnyContext, 11).Return(&order, nil)
				repoMock.UpdateOrderMock.Return(nil)

				prodMock.ProduceEventMock.Return(nil)

				cacheMock.GetMock.Expect(11).Return(&dto.OrderDTO{}, false)
				cacheMock.SetMock.Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ErrorOrderClientMismatch_GetRefundFromСlient",
			args: args{
				clientID: 10,
				orderID:  11,
			},
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				prodMock *mock.EventLogProducerFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {

				pickUpTime := time.Now()
				order := dto.OrderDTO{
					ID:         11,
					ClientID:   11,
					PickUpTime: sql.NullTime{Time: pickUpTime, Valid: true},
					Status:     domain.OrderStatusMap[domain.OrderStatusPickedUp],
				}

				repoMock.GetOrderByIDMock.Expect(minimock.AnyContext, 11).Return(&order, nil)

				cacheMock.GetMock.Expect(11).Return(&dto.OrderDTO{}, false)
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
			setup: func(
				repoMock *mock.OrderRepoFacadeMock,
				prodMock *mock.EventLogProducerFacadeMock,
				cacheMock *mock.OrderCacheFacadeMock,
			) {
				order := dto.OrderDTO{
					ID:       11,
					ClientID: 10,
					Status:   domain.OrderStatusMap[domain.OrderStatusReceived],
				}

				repoMock.GetOrderByIDMock.Expect(minimock.AnyContext, 11).Return(&order, nil)

				cacheMock.GetMock.Expect(11).Return(&dto.OrderDTO{}, false)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderIsNotRefundable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			repoMock := mock.NewOrderRepoFacadeMock(ctrl)
			prodMock := mock.NewEventLogProducerFacadeMock(ctrl)
			cacheMock := mock.NewOrderCacheFacadeMock(ctrl)

			tt.setup(repoMock, prodMock, cacheMock)
			uc := usecase.NewOrderUseCase(repoMock, prodMock, cacheMock)

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
		setup    func(*mock.OrderRepoFacadeMock)
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
			setup: func(repoMock *mock.OrderRepoFacadeMock) {
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

				repoMock.GetRefundsListMock.Expect(minimock.AnyContext, 0, 0).Return(refunds, nil)

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
			repoMock := mock.NewOrderRepoFacadeMock(ctrl)
			prodMock := mock.NewEventLogProducerFacadeMock(ctrl)
			cacheMock := mock.NewOrderCacheFacadeMock(ctrl)

			tt.setup(repoMock)
			uc := usecase.NewOrderUseCase(repoMock, prodMock, cacheMock)

			got, err := uc.RefundList(context.Background(), tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.RefundList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
