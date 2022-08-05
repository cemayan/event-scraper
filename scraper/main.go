package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/cemayan/event-scraper-common/protos"
	"github.com/cemayan/event-scraper/scraper/src/config"
	"github.com/cemayan/event-scraper/scraper/src/model"
	"github.com/cemayan/event-scraper/scraper/src/mq"
	"github.com/cemayan/event-scraper/scraper/src/scheduler"
	"github.com/cemayan/event-scraper/scraper/src/service"
	"github.com/cemayan/event-scraper/scraper/src/utils"
	"github.com/go-resty/resty/v2"
	"github.com/jasonlvhit/gocron"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var amqqConn *amqp.Connection
var amqqChannel *amqp.Channel
var grpcConn *grpc.ClientConn
var _log *logrus.Logger
var s *gocron.Scheduler
var configs = config.GetConfig()
var ampqService mq.MQClient
var biletixSvc service.ProviderService
var passoSvc service.ProviderService
var schedulerSvc scheduler.SchedulerService
var httpClient *resty.Client

func init() {

	//logrus init
	_log = logrus.New()
	_log.Out = os.Stdout

	str := fmt.Sprintf("%s:%s", configs.GRPC_ADDR, configs.GRPC_ADDR_PORT)
	//gRPC connection
	_grpcConn, err := grpc.Dial(str, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcConn = _grpcConn
	utils.FailOnError(err, "did not connect")

	_log.Infoln("gRPC connection is starting...")

	// RabbitMQ connection
	mqConn, err := amqp.Dial(configs.AMPQ_URI)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := mqConn.Channel()
	amqqChannel = ch
	utils.FailOnError(err, "Failed to open a channel")

	_log.Infoln("RabbitMQ connection is starting...")

	//Gocron init
	s = gocron.NewScheduler()

	//Resty http client
	httpClient = resty.New()

}

func main() {

	eventClient := pb.NewEventgRPCServiceClient(grpcConn)

	ampqService = mq.NewAMQPService(amqqChannel, _log)
	ampqService.ExchangeDeclare(configs.EXCHANGE_NAME, "direct", configs.QUEUE_NAME)

	biletixSvc = service.NewBiletixService(eventClient, _log, ampqService, httpClient)
	passoSvc = service.NewPassoService(eventClient, _log, ampqService, httpClient)
	schedulerSvc = scheduler.NewScheduler(s)

	if configs.SCHEDULE_ARRAY != "" {
		var scheduleModels []model.ScheduleModel
		json.Unmarshal([]byte(configs.SCHEDULE_ARRAY), &scheduleModels)

		for _, scheduleModel := range scheduleModels {

			if scheduleModel.Provider == utils.BILETIX {

				scheduler.SetJob(s,
					scheduleModel.TimeType,
					scheduleModel.TimeCount,
					scheduleModel.At,
					biletixSvc.GetEvents,
					utils.GetQueryParameters(scheduleModel.Category, scheduleModel.City))
			} else if scheduleModel.Provider == utils.PASSO {
				scheduler.SetJob(s,
					scheduleModel.TimeType,
					scheduleModel.TimeCount,
					scheduleModel.At,
					passoSvc.GetEvents,
					scheduleModel.City,
				)
			}

		}

		logrus.Infoln("gocron  is starting...")
		schedulerSvc.Start()

		defer grpcConn.Close()
		defer amqqConn.Close()
		defer amqqChannel.Close()

	}
}
