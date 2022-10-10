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
	"net/http"
	"time"
)

// BiletixSvc is representation of service dependencies
type BiletixSvc struct {
	grpcClient  pb.EventgRPCServiceClient
	log         *log.Logger
	ampqService mq.MQClient
	httpClient  *resty.Client
	configs     *scraper.AppConfig
}

// GetEvents scrapes the events based on given payload
// params includes the query strings
// Biletix query strings  example:
// start=0&rows=1300&q=*:*&fq=start%3A%5B2022-08-02T00%3A00%3A00Z%20TO%202022-09-01T00%3A00%3A00Z%2B1DAY%5D&sort=start%20asc,%20vote%20desc&&fq=category:%22MUSIC%22&fq=city:%22%C4%B0stanbul%22&wt=json&indent=true&facet=true&facet.field=category&facet.field=venuecode&facet.field=region&facet.field=subcategory&facet.mincount=1
// In order to send  events to DB , Grpc connection is started
// In addition, it must delete the all past events so "DELETE_EVENTS_IN_TABLE" event is sent
// This event is sent via amqp and includes which provider will be deleting
// Biletix prevents request without the cookie
// "BXID=AAAAAAVvcHP17piIDbV0DmuiSQXIxhBRoxUckbpoYxa/2QjFEQ==" cookie is used to send a http request
func (b BiletixSvc) GetEvents(params interface{}) {

	paramsArr := params.([]interface{})
	qs := paramsArr[0].(string)

	log.Infoln("Biletix event operation is starting...")

	stream, err := b.grpcClient.SendEvent(context.Background())

	var scraperEvent common.ScraperEvent

	payload := map[string]interface{}{}
	payload["provider"] = utils.BILETIX.String()

	scraperEvent.AggregationId = uuid.New()
	scraperEvent.EventDate = time.Now().Unix()
	scraperEvent.EventName = common.DELETE_EVENTS_IN_TABLE
	scraperEvent.Payload = payload

	eventBytes, err := json.Marshal(scraperEvent)
	utils.FailOnError(err, "event marshall error")

	b.ampqService.Publish(b.configs.RabbitMQ.EXCHANGE_NAME, b.configs.RabbitMQ.ROUTING_KEY, eventBytes)

	url := b.configs.Biletix.URL + qs

	b.httpClient.SetCookie(&http.Cookie{
		Name:     "BXID=AAAAAAVvcHP17piIDbV0DmuiSQXIxhBRoxUckbpoYxa/2QjFEQ==",
		Value:    "BXID=AAAAAAVvcHP17piIDbV0DmuiSQXIxhBRoxUckbpoYxa/2QjFEQ==",
		Path:     "/",
		Domain:   "biletix.com",
		MaxAge:   36000,
		HttpOnly: true,
		Secure:   false,
	})
	resp, err := b.httpClient.R().
		SetHeader("Content-Type", "application/json").
		EnableTrace().
		Get(url)

	utils.FailOnError(err, "While biletix event list getting received error ")

	var biletixResponse model.BiletixResponse
	err = json.Unmarshal(resp.Body(), &biletixResponse)
	utils.FailOnError(err, "Biletix json marshall error ")

	for _, doc := range biletixResponse.Response.Docs {

		var _event pb.Event

		_event.EventName = doc.Name
		_event.Place = doc.Venue
		_event.FirstDate = doc.Start.String()
		_event.SecondDate = doc.End.String()
		_event.Type = doc.Category
		_event.Provider = utils.BILETIX.String()

		if err := stream.Send(&_event); err != nil {
			log.Errorln(err)
		}
	}

	log.Infoln("Biletix  event operation is completed.")

	err = stream.CloseSend()
	utils.FailOnError(err, "While stream is closing there is a error")
}

func NewBiletixService(
	grpcClient pb.EventgRPCServiceClient,
	log *log.Logger,
	ampqService mq.MQClient,
	httpClient *resty.Client,
	configs *scraper.AppConfig,
) ProviderService {
	return &BiletixSvc{grpcClient: grpcClient, log: log, ampqService: ampqService, httpClient: httpClient, configs: configs}
}
