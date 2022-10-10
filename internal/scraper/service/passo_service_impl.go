package service

import (
	"context"
	"encoding/json"
	"github.com/cemayan/event-scraper/common"
	"github.com/cemayan/event-scraper/config/scraper"
	"github.com/cemayan/event-scraper/internal/scraper/model"
	"github.com/cemayan/event-scraper/internal/scraper/mq"
	"github.com/cemayan/event-scraper/internal/scraper/utils"
	pb "github.com/cemayan/event-scraper/protos"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

// PassoSvc is representation of service dependencies
type PassoSvc struct {
	grpcClient  pb.EventgRPCServiceClient
	log         *log.Logger
	ampqService mq.MQClient
	httpClient  *resty.Client
	configs     *scraper.AppConfig
}

// GetEvents scrapes the events based on given payload
// params includes the query strings
// In order to send  events to DB , Grpc connection is started
// In addition, it must delete the all past events so "DELETE_EVENTS_IN_TABLE" event is sent
// This event is sent via amqp and includes which provider will be deleting
// In order to send a http request to Passo it should be passed a payload
func (b PassoSvc) GetEvents(params interface{}) {

	log.Infoln("PassoResponse event operation is starting...")

	paramsArr := params.([]interface{})
	city := paramsArr[0].(utils.City)

	stream, err := b.grpcClient.SendEvent(context.Background())

	var scraperEvent common.ScraperEvent

	payload := map[string]interface{}{}
	payload["provider"] = utils.PASSO.String()

	scraperEvent.AggregationId = uuid.New()
	scraperEvent.EventDate = time.Now().Unix()
	scraperEvent.EventName = common.DELETE_EVENTS_IN_TABLE
	scraperEvent.Payload = payload

	eventBytes, err := json.Marshal(scraperEvent)
	utils.FailOnError(err, "event marshall error")

	b.ampqService.Publish(b.configs.RabbitMQ.EXCHANGE_NAME, b.configs.RabbitMQ.ROUTING_KEY, eventBytes)

	startDate, endDate := utils.GetDates()
	requestBody, err := json.Marshal(model.PassoRequestBody{
		CountRequired: true,
		HastagID:      nil,
		CityID:        city.String(),
		Date:          nil,
		VenueID:       nil,
		StartDate:     startDate,
		EndDate:       endDate,
		LanguageID:    100, //English
		From:          0,
		Size:          100,
	})

	utils.FailOnError(err, "marshall error!")

	resp, err := b.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		EnableTrace().
		Post(b.configs.Passo.URL)

	utils.FailOnError(err, "While passo event list getting it is go error ")

	var passoEvents model.PassoResponse
	err = json.Unmarshal(resp.Body(), &passoEvents)
	utils.FailOnError(err, "Passolig json marshall error ")

	for _, passoEvent := range passoEvents.ValueList {

		var _event pb.Event

		var _type string
		if len(passoEvent.HashTagList) != 0 {
			_type = passoEvent.HashTagList[0].HashTagName
		}

		_event.EventName = passoEvent.Name
		_event.Place = passoEvent.VenueName
		_event.FirstDate = passoEvent.Date
		_event.SecondDate = passoEvent.EndDate
		_event.Type = _type
		_event.Provider = utils.PASSO.String()

		if err := stream.Send(&_event); err != nil {
			log.Errorln(err)
		}
	}

	log.Infoln("PassoResponse event operation is completed.")

	err = stream.CloseSend()
	utils.FailOnError(err, "While stream is closing there is a error")

}

func NewPassoService(
	grpcClient pb.EventgRPCServiceClient,
	log *log.Logger,
	ampqService mq.MQClient,
	httpClient *resty.Client,
	configs *scraper.AppConfig,
) ProviderService {
	return &PassoSvc{grpcClient: grpcClient, log: log, ampqService: ampqService, httpClient: httpClient, configs: configs}
}
