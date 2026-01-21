package config

import (
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/exgamer/gosdk-rabbit-core/pkg/enums"
)

// NewPublisherFanoutDurableConfig
// Fanout: routing key не важен, topic можно игнорировать.
// Exchange фиксированный.
func NewPublisherFanoutDurableConfig(exchange string) []Config {
	return NewPublisherConfig(
		amqp.DefaultMarshaler{NotPersistentDeliveryMode: false}, // IMPORTANT: persistent messages
		amqp.PublishConfig{
			GenerateRoutingKey: func(topic string) string { return "" }, // fanout игнорирует
			ConfirmDelivery:    true,
		},
		amqp.ExchangeConfig{
			GenerateName: func(topic string) string { return exchange },
			Type:         enums.FanoutExchange,
			Durable:      true,
		},
		&amqp.DefaultTopologyBuilder{},
	)
}

// NewPublisherDirectDurableConfig
// Direct: topic = routing key
func NewPublisherDirectDurableConfig(exchange string) []Config {
	return NewPublisherConfig(
		amqp.DefaultMarshaler{NotPersistentDeliveryMode: false},
		amqp.PublishConfig{
			GenerateRoutingKey: func(topic string) string { return topic },
			ConfirmDelivery:    true,
		},
		amqp.ExchangeConfig{
			GenerateName: func(topic string) string { return exchange },
			Type:         enums.DirectExchange,
			Durable:      true,
		},
		&amqp.DefaultTopologyBuilder{},
	)
}

// NewPublisherTopicDurableConfig
// Topic exchange: topic = routing key (например "orders.paid")
func NewPublisherTopicDurableConfig(exchange string) []Config {
	return NewPublisherConfig(
		amqp.DefaultMarshaler{NotPersistentDeliveryMode: false},
		amqp.PublishConfig{
			GenerateRoutingKey: func(topic string) string { return topic },
			ConfirmDelivery:    true,
		},
		amqp.ExchangeConfig{
			GenerateName: func(topic string) string { return exchange },
			Type:         enums.TopicExchange,
			Durable:      true,
		},
		&amqp.DefaultTopologyBuilder{},
	)
}

// NewPublisherConfig - дефолтный конфиг для publisher (без биндинга очередей!)
func NewPublisherConfig(
	defaultMarshaller amqp.DefaultMarshaler,
	publishConfig amqp.PublishConfig,
	exchangeConfig amqp.ExchangeConfig,
	topologyBuilder *amqp.DefaultTopologyBuilder,
) []Config {
	return []Config{
		WithMarshaller(defaultMarshaller),
		WithPublishConfig(publishConfig),
		WithExchangeConfig(exchangeConfig),
		WithTopologyBuilder(topologyBuilder),
	}
}
