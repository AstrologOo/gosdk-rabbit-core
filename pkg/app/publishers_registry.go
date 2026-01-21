package app

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/exgamer/gosdk-rabbit-core/pkg/config"
	"github.com/exgamer/gosdk-rabbit-core/pkg/rabbitmq"
	"sync"
)

func NewPublisherRegistry(connection *amqp.ConnectionWrapper) *PublisherRegistry {
	return &PublisherRegistry{
		publishers:  make(map[string]*rabbitmq.Publisher),
		definitions: make(map[string]config.PublisherDefinition),
		connection:  connection,
	}
}

type PublisherRegistry struct {
	mu          sync.RWMutex
	publishers  map[string]*rabbitmq.Publisher
	definitions map[string]config.PublisherDefinition
	connection  *amqp.ConnectionWrapper
}

func (r *PublisherRegistry) RegisterMultiple(list []config.PublisherDefinition) error {
	for _, def := range list {
		if err := r.Register(def); err != nil {
			return err
		}
	}

	return nil
}

// Register — вызывается разработчиком в модуле (без conn!)
func (r *PublisherRegistry) Register(def config.PublisherDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if def.Name == "" {
		return fmt.Errorf("publisher name is empty")
	}

	if _, exists := r.definitions[def.Name]; exists {
		return fmt.Errorf("publisher %q already registered", def.Name)
	}

	r.definitions[def.Name] = def

	return nil
}

func (r *PublisherRegistry) Get(name string) (*rabbitmq.Publisher, error) {
	if name == "" {
		return nil, fmt.Errorf("publisher name is empty")
	}

	// Сначала пробуем быстро под RLock
	r.mu.RLock()
	pub, ok := r.publishers[name]
	r.mu.RUnlock()

	if ok {
		return pub, nil
	}

	// Если нет — создаём под Lock
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check (могли создать между RUnlock и Lock)
	if pub, ok = r.publishers[name]; ok {
		return pub, nil
	}

	def, exists := r.definitions[name]
	if !exists {
		return nil, fmt.Errorf("publisher %q is not registered", name)
	}

	if r.connection == nil {
		return nil, fmt.Errorf("rabbit connection is nil, cannot create publisher %q", name)
	}

	newPub, err := rabbitmq.NewAmqpPublisher(r.connection, def.Config...)
	if err != nil {
		return nil, fmt.Errorf("init publisher %q: %w", name, err)
	}

	r.publishers[name] = newPub

	return newPub, nil
}

func (r *PublisherRegistry) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var firstErr error

	for name, pub := range r.publishers {
		if pub == nil {
			continue
		}

		if err := pub.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close publisher %q: %w", name, err)
		}
	}

	// очищаем, чтобы registry можно было безопасно переинициализировать
	r.publishers = make(map[string]*rabbitmq.Publisher)

	return firstErr
}
