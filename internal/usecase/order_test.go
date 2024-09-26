package usecase_test

import (
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/domain"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/repository"
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
		setup    func(*mock.RepositoryMock)
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
			setup: func(repoMock *mock.RepositoryMock) {
				order, _ := domain.NewOrder(dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
					Packages:   []string{"unknown", "unknown"},
				}, nil, nil)

				repoMock.AddOrderMock.Expect(order).Return(nil)
				repoMock.UpdateMock.Expect().Return(nil)
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
			setup: func(repoMock *mock.RepositoryMock) {
				order, _ := domain.NewOrder(dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
					Packages:   []string{"unknown", "unknown"},
				}, nil, nil)
				repoMock.AddOrderMock.Expect(order).Return(repository.ErrAlreadyExist)
			},
			wantErr:  true,
			errValue: repository.ErrAlreadyExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			repoMock := mock.NewRepositoryMock(ctrl)

			tt.setup(repoMock)

			uc := usecase.NewOrderUseCase(repoMock)

			err := uc.ReceiveOrderFromCourier(tt.args.req)
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
		setup    func(*mock.RepositoryMock)
		wantErr  bool
		errValue error
	}{
		{
			name: "Success_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(repoMock *mock.RepositoryMock) {
				var order domain.Order
				order.SetClientID(10)
				order.SetStoreUntil(time.Now())
				order.SetStatus(domain.OrderStatusReceived)

				repoMock.GetOrderByIDMock.Expect(10).Return(&order, nil)
				repoMock.UpdateMock.Expect().Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ErrorOrderPickedUp_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(repoMock *mock.RepositoryMock) {
				var order domain.Order
				order.SetClientID(10)
				order.SetStoreUntil(time.Now())
				order.SetStatus(domain.OrderStatusPickedUp)

				repoMock.GetOrderByIDMock.Expect(10).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderPickedUp,
		},
		{
			name: "ErrorOrderDeleted_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(repoMock *mock.RepositoryMock) {
				var order domain.Order
				order.SetClientID(10)
				order.SetStoreUntil(time.Now())
				order.SetStatus(domain.OrderStatusDelete)

				repoMock.GetOrderByIDMock.Expect(10).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderDeleted,
		},
		{
			name: "ErrorOrderStoreTimeNotExpired_ReturnOrderToCourier",
			args: args{orderID: 10},
			setup: func(repoMock *mock.RepositoryMock) {
				var order domain.Order
				order.SetClientID(10)
				order.SetStoreUntil(time.Now().Add(24 * time.Hour))
				order.SetStatus(domain.OrderStatusReceived)

				repoMock.GetOrderByIDMock.Expect(10).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderStoreTimeNotExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			repoMock := mock.NewRepositoryMock(ctrl)
			tt.setup(repoMock)
			uc := usecase.NewOrderUseCase(repoMock)

			err := uc.ReturnOrderToCourier(tt.args.orderID)
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
		orderIDs []int
	}

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.RepositoryMock)
		wantErr  bool
		errValue error
	}{
		{
			name: "SuccessReturnOrderToCourier",
			args: args{orderIDs: []int{10}},
			setup: func(repoMock *mock.RepositoryMock) {
				var orders []*domain.Order

				var order domain.Order
				order.SetClientID(10)
				order.SetStoreUntil(time.Now().Add(24 * time.Hour))
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

			ctrl := minimock.NewController(t)
			repoMock := mock.NewRepositoryMock(ctrl)
			tt.setup(repoMock)
			uc := usecase.NewOrderUseCase(repoMock)

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

	successStoreTime := time.Now().Add(48 * time.Hour)

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.RepositoryMock)
		want     *dto.ListOrdersDTO
		wantErr  bool
		errValue error
	}{
		{
			name: "SuccessOrderList",
			args: args{
				clientID: 10,
			},
			setup: func(repoMock *mock.RepositoryMock) {
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
			t.Parallel()

			ctrl := minimock.NewController(t)
			repoMock := mock.NewRepositoryMock(ctrl)
			tt.setup(repoMock)
			uc := usecase.NewOrderUseCase(repoMock)

			got, err := uc.OrderList(tt.args.clientID)
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
		orderID  int
	}

	tests := []struct {
		name     string
		args     args
		setup    func(*mock.RepositoryMock)
		wantErr  bool
		errValue error
	}{
		{
			name: "Success_GetRefundFromСlient",
			args: args{
				clientID: 10,
				orderID:  11,
			},
			setup: func(repoMock *mock.RepositoryMock) {
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
		{
			name: "ErrorOrderClientMismatch_GetRefundFromСlient",
			args: args{
				clientID: 10,
				orderID:  11,
			},
			setup: func(repoMock *mock.RepositoryMock) {
				var order domain.Order
				order.SetID(11)
				order.SetClientID(11)
				order.SetPickUpTime()
				order.SetStatus(domain.OrderStatusPickedUp)

				repoMock.GetOrderByIDMock.Expect(11).Return(&order, nil)
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
			setup: func(repoMock *mock.RepositoryMock) {
				var order domain.Order
				order.SetID(11)
				order.SetClientID(10)
				order.SetStatus(domain.OrderStatusReceived)

				repoMock.GetOrderByIDMock.Expect(11).Return(&order, nil)
			},
			wantErr:  true,
			errValue: usecase.ErrOrderIsNotRefundable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			repoMock := mock.NewRepositoryMock(ctrl)
			tt.setup(repoMock)
			uc := usecase.NewOrderUseCase(repoMock)

			err := uc.GetRefundFromСlient(tt.args.clientID, tt.args.orderID)

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
		setup    func(*mock.RepositoryMock)
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
			setup: func(repoMock *mock.RepositoryMock) {
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
			t.Parallel()

			ctrl := minimock.NewController(t)
			repoMock := mock.NewRepositoryMock(ctrl)
			tt.setup(repoMock)
			uc := usecase.NewOrderUseCase(repoMock)

			got, err := uc.RefundList(tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderUseCase.RefundList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
