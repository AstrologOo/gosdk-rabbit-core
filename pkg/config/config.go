package config

import (
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
)

type Config func(*amqp.Config)

// WithAmqpURI - конфиг лоя лобавления ссылки коннект
func WithAmqpURI(uri string) Config {
	return func(config *amqp.Config) {
		config.Connection.AmqpURI = uri
	}
}

// WithMarshaller - конфиг для маршалинга сообщений
func WithMarshaller(m amqp.DefaultMarshaler) Config {
	return func(config *amqp.Config) {
		config.Marshaler = m
	}
}

// WithExchangeConfig - add exchange ExchangeConfig
func WithExchangeConfig(ec amqp.ExchangeConfig) Config {
	return func(config *amqp.Config) {
		config.Exchange = ec
	}
}

// WithPublishConfig - add publish config PublishConfig
func WithPublishConfig(pc amqp.PublishConfig) Config {
	return func(config *amqp.Config) {
		config.Publish = pc
	}
}

// WithConsumerConfig - для добавление конфига консьюмера
func WithConsumerConfig(cc amqp.ConsumeConfig) Config {
	return func(config *amqp.Config) {
		config.Consume = cc
	}
}

// WithRoutingKeyBinding - добавление конфига для бинда с routing key
func WithRoutingKeyBinding(rk amqp.QueueBindConfig) Config {
	return func(config *amqp.Config) {
		config.QueueBind = rk
	}
}

// WithQueueName  - добавление конфгиа для названия очереди
func WithQueueName(q amqp.QueueConfig) Config {
	return func(config *amqp.Config) {
		config.Queue = q
	}
}

// WithTopologyBuilder
func WithTopologyBuilder(tp amqp.TopologyBuilder) Config {
	if tp == nil {
		tp = &amqp.DefaultTopologyBuilder{}
	}

	return func(config *amqp.Config) {
		config.TopologyBuilder = tp
	}
}
