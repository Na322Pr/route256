package producer

import (
	"github.com/IBM/sarama"
)

func PrepareConfig(opts ...Option) *sarama.Config {
	с := sarama.NewConfig()

	с.Producer.Return.Successes = true
	с.Producer.Return.Errors = true

	for _, opt := range opts {
		_ = opt.Apply(с)
	}

	return с
}
