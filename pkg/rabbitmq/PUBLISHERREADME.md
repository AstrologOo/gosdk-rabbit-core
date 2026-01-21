# RabbitMQ Publisher (Watermill AMQP)

Этот пакет предоставляет **долгоживущий** RabbitMQ publisher поверх **общего AMQP connection** (Watermill AMQP), а также **реестр паблишеров**, чтобы модули могли удобно объявлять свои паблишеры, а `RabbitKernel` — инициализировать их на старте.

---

## Зачем так сделано

### ✅ Правильный lifecycle
- AMQP connection создаётся **один раз** в `RabbitKernel`.
- Каждый publisher создаётся **один раз** и переиспользуется на протяжении жизни приложения.
- На `Stop()` все publisher’ы корректно закрываются, затем закрывается connection.

### ❌ Неправильный вариант
Создавать `amqp.NewPublisherWithConnection(...)` на каждый `Publish()` — дорого и ломает смысл channel pool и confirm delivery.

---

## Термины: что такое `topic`

В Watermill `topic` — это логическое имя канала сообщений.  
В нашем SDK принято соглашение:

- Для `direct` и `topic` exchange:  
  **`topic` = routing key**
- Для `fanout` exchange:  
  routing key игнорируется

Пример:
```go
events.Publish("orders.paid", payload)
```

RabbitMQ:
```
exchange = "events"
routing_key = "orders.paid"
```

---

## Быстрый старт

Включить publisher в Kernel:

```go
rabbitKernel := app.NewRabbitKernel().
    EnablePublisher()
```

---

## Объявление паблишеров

```go
pubReg.Register(config.NewPublisherDefinition(
    "events",
    config.NewPublisherTopicDurableConfig("events")...,
))

pubReg.Register(config.NewPublisherDefinition(
    "broadcast",
    config.NewPublisherFanoutDurableConfig("broadcast")...,
))
```

- `events`, `broadcast` — имена паблишеров в registry
- `"events"`, `"broadcast"` — RabbitMQ exchange

---

## Использование

```go
eventsPub, err := pubReg.Get("events")
if err != nil {
    return err
}

return eventsPub.Publish("orders.paid", payload)
```

С метаданными:

```go
eventsPub.PublishWithMetaData("orders.paid", map[string]string{
    "trace_id": traceID,
}, payload)
```

Batch:

```go
eventsPub.PublishBatch("products.updated", []any{p1, p2, p3})
```

---

## Durable delivery

`*DurableConfig` означает:
- Exchange с `Durable=true`
- Сообщения публикуются как **persistent**

Это гарантирует сохранность сообщений при рестарте RabbitMQ.

---

## API

### Publisher
- `Publish(topic string, payload any) error`
- `PublishWithMetaData(topic string, meta map[string]string, payload any) error`
- `PublishBatch(topic string, payloads []any) error`
- `PublishBatchWithMetaData(topic string, meta map[string]string, payloads []any) error`
- `Close() error`

### PublisherRegistry
- `Register(def PublisherDefinition) error`
- `RegisterMultiple(list []PublisherDefinition) error`
- `Init(conn *amqp.ConnectionWrapper) error`
- `Get(name string) (*Publisher, error)`
- `Close() error`

---

## Рекомендации по неймингу routing keys

Используйте формат:

```
domain.event
```

Примеры:
- `orders.paid`
- `orders.created`
- `products.updated`
- `users.created`

Позволяет подписываться масками:
- `orders.*`
- `*.updated`
