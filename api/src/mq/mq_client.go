package mq

import amqp "github.com/rabbitmq/amqp091-go"

// MQClient is representation of a RabbitMQ client operations
type MQClient interface {
	ExchangeDeclare(exchangeName string, exchangeType string, queueName string)
	QueueBind(queueName string, routingKey string, exchangeName string)
	Consumer(queueName string, consumerTag string, autoAck bool) <-chan amqp.Delivery
	HandleConsumer(deliveries <-chan amqp.Delivery)
	Publish(exchangeName string, routingKey string, body []byte)
}
