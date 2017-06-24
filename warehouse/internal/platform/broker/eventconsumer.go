package broker

import (
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/broker"
	"github.com/micro/protobuf/proto"
)

const (
	topic = "go.shopping.item.shipped"
)

// CreateEventConsumer creates a broker subscription that converts broker messages into
// item shipped events, placing those events on that channel so that they can be processed
// by other modules.
func CreateEventConsumer(itemShippedChannel chan *shipping.ItemShippedEvent) (err error) {
	_, err = broker.Subscribe(topic, func(p broker.Publication) error {
		log.Logf("[sub] received message %+v", p.Message().Header)

		var shippedEvent shipping.ItemShippedEvent
		err = proto.Unmarshal(p.Message().Body, &shippedEvent)
		if err != nil {
			log.Logf("Failed to unmarshal broker message: %s", err)
			return err
		}
		itemShippedChannel <- &shippedEvent
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
