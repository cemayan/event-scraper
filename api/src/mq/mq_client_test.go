package mq

import (
	"encoding/json"
	"github.com/cemayan/event-scraper-common/events"
	"github.com/cemayan/event-scraper-common/protos"
	"github.com/cemayan/event-scraper/api/src/config"
	"github.com/cemayan/event-scraper/api/src/database"
	"github.com/cemayan/event-scraper/api/src/repo"
	"github.com/cemayan/event-scraper/api/src/service"
	"github.com/cemayan/event-scraper/api/src/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

var configs = config.GetConfig()

type e2eTestSuite struct {
	suite.Suite
	app      *fiber.App
	db       *gorm.DB
	mqClient MQClient
	channel  *amqp.Channel
	eSvc     service.EventService
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (ts *e2eTestSuite) removeAllRecords() {
	ts.db.Exec("DELETE FROM events")
}

func (ts *e2eTestSuite) createSomeRecord() {
	eventModel := protos.Event{

		Type:       "MUSIC",
		EventName:  "TEST_EVENT",
		Place:      "TEST_PLACE",
		FirstDate:  "2022-08-17 18:00:00 +0000 UTC",
		SecondDate: "2022-08-17 18:00:00 +0000 UTC",
		Provider:   "BILETIX",
	}

	eventModel2 := protos.Event{
		Type:       "ART",
		EventName:  "TEST_EVENT2",
		Place:      "TEST_PLACE2",
		FirstDate:  "2022-08-17 18:00:00 +0000 UTC",
		SecondDate: "2022-08-17 18:00:00 +0000 UTC",
		Provider:   "BILETIX",
	}

	eventModel3 := protos.Event{
		Type:       "MUSIC",
		EventName:  "TEST_EVENT3",
		Place:      "TEST_PLACE3",
		FirstDate:  "2022-08-17 18:00:00 +0000 UTC",
		SecondDate: "2022-08-17 18:00:00 +0000 UTC",
		Provider:   "PASSO",
	}

	ts.db.Create(&eventModel)
	ts.db.Create(&eventModel2)
	ts.db.Create(&eventModel3)
}

func (ts *e2eTestSuite) getRecords() []protos.Event {
	var events []protos.Event
	ts.db.Find(&events)
	return events
}

func (ts *e2eTestSuite) SetupSuite() {

	app := fiber.New()
	app.Use(cors.New())

	ts.app = app
	DB := database.GetDB()
	ts.db = DB

	// RabbitMQ connection
	mqConn, err := amqp.Dial(configs.AMPQ_URI)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := mqConn.Channel()
	ts.channel = ch
	utils.FailOnError(err, "Failed to open a channel")

	eRepo := repo.NewEventRepo(DB, log.New())
	eSvc := service.NewEventService(eRepo, log.New())
	ts.eSvc = eSvc

	amqpSvc := NewAMQPService(ch, eSvc, log.New())
	ts.mqClient = amqpSvc

	ts.mqClient.ExchangeDeclare(configs.EXCHANGE_NAME, "direct", configs.QUEUE_NAME)
	ts.mqClient.QueueBind(configs.QUEUE_NAME, configs.ROUTING_KEY, configs.EXCHANGE_NAME)

	consumer := ts.mqClient.Consumer(configs.QUEUE_NAME, configs.CONSUMER_TAG, true)
	go ts.mqClient.HandleConsumer(consumer)
}

func (ts *e2eTestSuite) TestEventService_HandleConsumer() {

	ts.removeAllRecords()
	ts.createSomeRecord()

	payload := map[string]interface{}{}
	payload["provider"] = utils.BILETIX.String()

	var event events.ScraperEvent
	event.EventDate = 1659469217
	event.EventName = events.DELETE_EVENTS_IN_TABLE
	event.Payload = payload
	event.AggregationId = uuid.New()

	bytesEvent, err := json.Marshal(event)
	if err != nil {
		return
	}

	// This event will remove based on given provider
	ts.mqClient.Publish(configs.EXCHANGE_NAME, configs.ROUTING_KEY, bytesEvent)
	time.Sleep(10 * time.Second)
	_events := ts.getRecords()
	ts.Equal(1, len(_events))
}
