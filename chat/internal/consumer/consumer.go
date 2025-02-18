package consumer

import (
	"encoding/json"
	"github.com/HJyup/translatify-chat/internal/models"
	"github.com/HJyup/translatify-common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	service models.ChatService
}

func NewConsumer(service models.ChatService) *Consumer {
	return &Consumer{service: service}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.MessageTranslatedEvent, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	var forever chan struct{}

	go func() {
		for d := range msg {

			msg := &models.ConsumerResponse{}
			err = json.Unmarshal(d.Body, msg)
			if err != nil {
				panic(err)
			}

			err = c.service.UpdateMessageTranslation(msg.MessageId, msg.TranslatedContent)
			if err != nil {
				panic(err)
			}

			log.Println("Message is updated")
		}
	}()

	<-forever
}
