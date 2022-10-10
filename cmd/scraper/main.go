package main

import (
	"encoding/json"
	"fmt"
	"github.com/cemayan/event-scraper/config/scraper"
	"github.com/cemayan/event-scraper/internal/scraper/model"
	"github.com/cemayan/event-scraper/internal/scraper/mq"
	"github.com/cemayan/event-scraper/internal/scraper/scheduler"
	"github.com/cemayan/event-scraper/internal/scraper/service"
	"github.com/cemayan/event-scraper/internal/scraper/utils"
	pb "github.com/cemayan/event-scraper/protos"
	"github.com/go-resty/resty/v2"
	"github.com/jasonlvhit/gocron"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var amqqConn *amqp.Connection
var amqqChannel *amqp.Channel
var grpcConn *grpc.ClientConn
var _log *logrus.Logger
var s *gocron.Scheduler
var ampqService mq.MQClient
var biletixSvc service.ProviderService
var passoSvc service.ProviderService
var schedulerSvc scheduler.SchedulerService
var httpClient *resty.Client
var configs *scraper.AppConfig
var v *viper.Viper

func init() {

	//logrus init
	_log = logrus.New()
	_log.Out = os.Stdout

	v = viper.New()
	_configs := scraper.NewConfig(v)

	env := os.Getenv("ENV")
	appConfig, err := _configs.GetConfig(env)
	configs = appConfig
	if err != nil {
		return
	}

	str := fmt.Sprintf("%s:%s", configs.Grpc.ADDR, configs.Grpc.PORT)
	//gRPC connection
	_grpcConn, err := grpc.Dial(str, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcConn = _grpcConn
	utils.FailOnError(err, "did not connect")

	_log.Infoln("gRPC connection is starting...")

	// RabbitMQ connection
	mqConn, err := amqp.Dial(configs.RabbitMQ.AMPQ_URI)
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
	ampqService.ExchangeDeclare(configs.RabbitMQ.EXCHANGE_NAME, "direct", configs.RabbitMQ.QUEUE_NAME)

	biletixSvc = service.NewBiletixService(eventClient, _log, ampqService, httpClient, configs)
	passoSvc = service.NewPassoService(eventClient, _log, ampqService, httpClient, configs)
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
