package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

func TestNewOrder(t *testing.T) {
	type args struct {
		orderDTO    dto.AddOrder
		packOpts    []PackageOption
		domainError error
	}

	successStoreTime := time.Now().Add(48 * time.Hour)

	tests := []struct {
		name    string
		args    args
		want    *Order
		wantErr bool
	}{
		{
			name: "SuccessNewOrder",
			args: args{
				orderDTO: dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
				},
				domainError: nil,
			},
			want: &Order{
				id:         1,
				clientID:   1,
				storeUntil: successStoreTime,
				status:     OrderStatusReceived,
				cost:       1000,
				weight:     5,
			},
			wantErr: false,
		},
		{
			name: "SuccessWithPackageNewOrder",
			args: args{
				orderDTO: dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
					Packages:   []string{"bag", "tape"},
				},
				domainError: nil,
			},
			want: &Order{
				id:         1,
				clientID:   1,
				storeUntil: successStoreTime,
				status:     OrderStatusReceived,
				cost:       1006,
				weight:     5,
				packages:   []OrderPackage{OrderPackageBag, OrderPackageTape},
			},
			wantErr: false,
		},
		{
			name: "ErrorInvalidIDNewOrder",
			args: args{
				orderDTO: dto.AddOrder{
					ID:         -1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
				},
				domainError: ErrInvalidID,
			},
			wantErr: true,
		},
		{
			name: "ErrorInvalidClientIDNewOrder",
			args: args{
				orderDTO: dto.AddOrder{
					ID:         1,
					ClientID:   -1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
				},
				domainError: ErrInvalidClientID,
			},
			wantErr: true,
		},
		{
			name: "ErrorInvalidCostNewOrder",
			args: args{
				orderDTO: dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       -120,
					Weight:     5,
				},
				domainError: ErrInvalidCost,
			},
			wantErr: true,
		},
		{
			name: "ErrorInvalidWeightNewOrder",
			args: args{
				orderDTO: dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     -15,
				},
				domainError: ErrInvalidWeight,
			},
			wantErr: true,
		},
		{
			name: "ErrorStoreTimeExpiredNewOrder",
			args: args{
				orderDTO: dto.AddOrder{
					ID:         1,
					ClientID:   1,
					StoreUntil: time.Now(),
					Cost:       1000,
					Weight:     5,
				},
				packOpts:    []PackageOption{},
				domainError: ErrStoreTimeExpired,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewOrder(tt.args.orderDTO)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.args.domainError)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
