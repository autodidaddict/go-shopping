package broker

import (
	"fmt"
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"github.com/micro/go-micro/broker"
	"github.com/micro/protobuf/proto"
	"log"
)

const (
	topic = "go.shopping.item.shipped"
)

// EventPublisher is an event publisher for the go-micro broker
type EventPublisher struct{}

// NewEventPublisher creates a new broker event publisher
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{}
}

// PublishItemShippedEvent publishes an item shipped event on the broker
func (p *EventPublisher) PublishItemShippedEvent(event *shipping.ItemShippedEvent) (err error) {
	bytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	msg := &broker.Message{
		Header: map[string]string{
			"sku":      event.Sku,
			"order-id": fmt.Sprintf("%d", event.OrderId),
		},
		Body: bytes,
	}
	if err := broker.Publish(topic, msg); err != nil {
		log.Printf("[pub] failed: %v\n", err)
		return err
	}
	log.Printf("[pub] pubbed item shipped event, %s/%d", event.Sku, event.OrderId)
	return nil
}
