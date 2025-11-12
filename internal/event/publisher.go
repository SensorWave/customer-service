package event

import (
    "encoding/json"
    "github.com/streadway/amqp"
    "log"
)

type Publisher struct {
    Channel  *amqp.Channel
    Exchange string
}

func NewPublisher(ch *amqp.Channel, exchange string) *Publisher {
    return &Publisher{Channel: ch, Exchange: exchange}
}

func (p *Publisher) Publish(routingKey string, payload interface{}) {
    body, _ := json.Marshal(payload)
    if err := p.Channel.Publish(
        p.Exchange,
        routingKey,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    ); err != nil {
        log.Println("Failed to publish event:", err)
    }
}
