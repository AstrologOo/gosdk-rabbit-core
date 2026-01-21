package rabbitmq

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/exgamer/gosdk-core/pkg/logger"
)

func NewAmqpConnection(cfg amqp.ConnectionConfig) (*amqp.ConnectionWrapper, error) {
	debug := false
	trace := false
	if logger.IsDebugLevel() {
		debug = true
	}

	if logger.GetLevel() <= logger.LevelTrace {
		trace = true
	}

	conn, err := amqp.NewConnection(cfg, watermill.NewStdLogger(debug, trace))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
