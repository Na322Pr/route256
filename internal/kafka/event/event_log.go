package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/marchenkosasha2/homework/internal/dto"
)

type EventType string

const (
	EventTypeReceive EventType = "receive"
	EventTypeGiveOut EventType = "giveout"
	EventTypeRefund  EventType = "refund"
)

type Event struct {
	Order           dto.OrderDTO `json:"order_info"`
	EventType       string       `json:"event"`
	OperationMoment time.Time    `json:"moment"`
}

type ProdFacade interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type EventLogProducer struct {
	prod    ProdFacade
	topic   string
	appName string
}

func NewEventLogProducer(prod ProdFacade, topic, appName string) (*EventLogProducer, error) {
	return &EventLogProducer{
		prod:    prod,
		topic:   topic,
		appName: appName,
	}, nil
}

func (ep *EventLogProducer) ProduceEvent(order dto.OrderDTO, eventType EventType) error {
	op := "EventLogProducer.ProduceEvent"

	event := &Event{
		Order:           order,
		EventType:       string(eventType),
		OperationMoment: time.Now(),
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	msg := &sarama.ProducerMessage{
		Topic: ep.topic,
		Value: sarama.ByteEncoder(bytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("app-name"),
				Value: []byte(ep.appName),
			},
		},
		Timestamp: time.Now(),
	}

	_, _, err = ep.prod.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
