package consumer

import (
	"encoding/json"
	"github.com/HJyup/translatify-common/broker"
	"github.com/HJyup/translatify-translation/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	service models.TranslationService
}

func NewConsumer(service models.TranslationService) *Consumer {
	return &Consumer{service: service}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.MessageSentEvent, true, false, false, false, nil)
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

			trnResponse, err := c.service.TranslateMessage(msg.Content, msg.SourceLang, msg.TargetLang)
			if err != nil {
				panic(err)
			}

			log.Println(trnResponse.TranslatedContent)
		}
	}()

	<-forever
}
