package main

import (
	"github.com/cemayan/event-scraper/api/src/config"
	"github.com/cemayan/event-scraper/api/src/database"
	"github.com/cemayan/event-scraper/api/src/mq"
	"github.com/cemayan/event-scraper/api/src/repo"
	"github.com/cemayan/event-scraper/api/src/service"
	"github.com/cemayan/event-scraper/api/src/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"os"
)

var _log *logrus.Logger
var eventRepo repo.EventRepository
var amqpChannel *amqp.Channel
var ampqService mq.MQClient
var eventService service.EventService
var configs = config.GetConfig()

func init() {
	//logrus init
	_log = logrus.New()
	_log.Out = os.Stdout

	// DB connection
	database.ConnectDB(configs)
	eventRepo = repo.NewEventRepo(database.GetDB(), _log)

	// RabbitMQ connection
	mqConn, err := amqp.Dial(configs.AMPQ_URI)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := mqConn.Channel()
	amqpChannel = ch
	utils.FailOnError(err, "Failed to open a channel")

	_log.Infoln("RabbitMQ connection is starting...")
}

func main() {
	start := make(chan bool)

	eventService = service.NewEventService(eventRepo, _log)
	ampqService = mq.NewAMQPService(amqpChannel, eventService, _log)

	// Exchange and queue is created
	ampqService.ExchangeDeclare(configs.EXCHANGE_NAME, "direct", configs.QUEUE_NAME)

	ampqService.QueueBind(configs.QUEUE_NAME, configs.ROUTING_KEY, configs.EXCHANGE_NAME)
	consumer := ampqService.Consumer(configs.QUEUE_NAME, configs.CONSUMER_TAG, true)
	go ampqService.HandleConsumer(consumer)

	<-start
}
