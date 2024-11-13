package event_test

import (
	"testing"
	"time"

	"github.com/Na322Pr/route256/internal/dto"
	"github.com/Na322Pr/route256/internal/kafka/event"
	"github.com/Na322Pr/route256/internal/kafka/event/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestEventLogProducer_ProduceEvent(t *testing.T) {
	type args struct {
		reqOrder     dto.OrderDTO
		reqEventType event.EventType
	}

	successStoreTime := time.Now().Add(48 * time.Hour)

	tests := []struct {
		name  string
		args  args
		setup func(
			*mock.ProdFacadeMock,
		)
		wantErr  bool
		errValue error
	}{
		{
			name: "SuccessProduceEvent",
			args: args{
				reqOrder: dto.OrderDTO{
					ID:         1,
					ClientID:   1,
					StoreUntil: successStoreTime,
					Cost:       1000,
					Weight:     5,
				},
				reqEventType: event.EventTypeReceive,
			},
			setup: func(
				prodMock *mock.ProdFacadeMock,
			) {
				prodMock.SendMessageMock.Return(0, 0, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			prodMock := mock.NewProdFacadeMock(ctrl)

			tt.setup(prodMock)

			ep, _ := event.NewEventLogProducer(prodMock, "pvz.events-log", "pvz-service")

			err := ep.ProduceEvent(tt.args.reqOrder, tt.args.reqEventType)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.errValue)
				return
			}
		})
	}
}
