package producer

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Na322Pr/route256/internal/config"
)

func NewSyncProducer(cfg config.Kafka, opts ...Option) (sarama.SyncProducer, error) {

	config := PrepareConfig(opts...)

	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return syncProducer, nil
}
