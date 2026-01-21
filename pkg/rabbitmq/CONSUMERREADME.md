## 🧩 Работа с консьюмером

Консьюмеры регистрируются в DI и доступен в бизнес-модулях.

```go
package consumers

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewCityConsumer() *CityConsumer {
	return &CityConsumer{}
}

type CityConsumer struct {
}

func (c *CityConsumer) Consume(ctx context.Context, msg *message.Message) error {
	fmt.Println(string(msg.Payload))

	return nil
}
```

```go
package city

import (
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/internal/domains/handbbok/modules/city/factories"
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/pkg/config"
)

func GetConsumers(
	consumersFactory *factories.CityConsumersFactory,
) []config.HandlerRegister {
	return []config.HandlerRegister{
		{
			Handler: consumersFactory.CityConsumer.Consume,
			Config: config.NewConsumerTopicDurableConfig(
				"test-rk",
				"test",
				"test",
				100,
			),
		},
	}
}

```

```go
consumersFactory := factories.NewCityConsumersFactory()

consumers := GetConsumers(consumersFactory)

reg, err := di.GetRabbitConsumersRegistry(a.Container) // твой helper для DI
if err != nil {
    return err
}

reg.RegisterMultipleHandler(consumers)
```

---