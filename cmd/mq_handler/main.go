package main

import (
	"github.com/cemayan/event-scraper/config/api"
	"github.com/cemayan/event-scraper/internal/api/database"
	"github.com/cemayan/event-scraper/internal/api/mq"
	"github.com/cemayan/event-scraper/internal/api/repo"
	"github.com/cemayan/event-scraper/internal/api/service"
	"github.com/cemayan/event-scraper/internal/api/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var _log *logrus.Logger
var eventRepo repo.EventRepository
var amqpChannel *amqp.Channel
var ampqService mq.MQClient
var eventService service.EventService
var configs *api.AppConfig
var v *viper.Viper
var dbHandler database.DBHandler

func init() {
	//logrus init
	_log = logrus.New()
	_log.Out = os.Stdout

	// DB connection
	v = viper.New()
	_configs := api.NewConfig(v)

	env := os.Getenv("ENV")
	appConfig, err := _configs.GetConfig(env)
	configs = appConfig
	if err != nil {
		return
	}

	//Postresql connection
	dbHandler = database.NewDbHandler(configs)
	dbHandler.ConnectDB()

	eventRepo = repo.NewEventRepo(database.DB, _log)

	// RabbitMQ connection
	mqConn, err := amqp.Dial(configs.RabbitMQ.AMPQ_URI)
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
	ampqService.ExchangeDeclare(configs.RabbitMQ.EXCHANGE_NAME, "direct", configs.RabbitMQ.QUEUE_NAME)

	ampqService.QueueBind(configs.RabbitMQ.QUEUE_NAME, configs.RabbitMQ.ROUTING_KEY, configs.RabbitMQ.EXCHANGE_NAME)
	consumer := ampqService.Consumer(configs.RabbitMQ.QUEUE_NAME, configs.RabbitMQ.CONSUMER_TAG, true)
	go ampqService.HandleConsumer(consumer)

	<-start
}
