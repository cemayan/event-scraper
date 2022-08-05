package mq

import (
	"context"
	"encoding/json"
	"github.com/cemayan/event-scraper-common/events"
	"github.com/cemayan/event-scraper/api/src/service"
	"github.com/cemayan/event-scraper/api/src/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

// MQCli is representation of a dependencies
type MQCli struct {
	channel  *amqp.Channel
	eventSvc service.EventService
	log      *log.Logger
}

// Publish servers to publish messages on RabbitMQ
func (M MQCli) Publish(exchangeName string, routingKey string, body []byte) {
	//seqNo := M.channel.GetNextPublishSeqNo()
	M.log.Printf("publishing %dB body (%q)", len(body), body)

	if err := M.channel.PublishWithContext(
		context.Background(),
		exchangeName, // publish to an exchange
		routingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(body),
			DeliveryMode: amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:     0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return
	}
}

// HandleConsumer consumes incoming events
// Based on given event provider it is deleted all provider events
func (M MQCli) HandleConsumer(deliveries <-chan amqp.Delivery) {
	M.log.Infoln("Channel consume operation is starting...")
	for tt := range deliveries {
		var event events.ScraperEvent
		err := json.Unmarshal(tt.Body, &event)
		if err != nil {
			continue
		}

		if event.EventName == events.DELETE_EVENTS_IN_TABLE {
			payload := event.Payload.(map[string]interface{})
			M.eventSvc.DeleteByProvider(payload["provider"].(string))
		}
	}

	M.log.Infoln("Channel consume operation is completed...")

}

// Consumer returns a channel based on given queueName, consumerTag and autoAck
func (M MQCli) Consumer(queueName string, consumerTag string, autoAck bool) <-chan amqp.Delivery {
	deliveries, err := M.channel.Consume(
		queueName,   // name
		consumerTag, // consumerTag,
		autoAck,     // autoAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // arguments
	)

	utils.FailOnError(err, "Queue Consume")

	return deliveries
}

// QueueBind is used to bind a  queue
func (M MQCli) QueueBind(queueName string, routingKey string, exchangeName string) {

	err := M.channel.QueueBind(
		queueName,    // name of the queue
		routingKey,   // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)
	utils.FailOnError(err, "Queue Bind")
}

// ExchangeDeclare is used to declare a  queue and exchange
func (M MQCli) ExchangeDeclare(exchangeName string, exchangeType string, queueName string) {

	_, err := M.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	err = M.channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	utils.FailOnError(err, "Exchange Declare")
}

func NewAMQPService(channel *amqp.Channel, eventSvc service.EventService, log *log.Logger) MQClient {
	return &MQCli{channel: channel, eventSvc: eventSvc, log: log}
}
