package test

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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"os"
	"testing"
)

type e2eScraperTestSuite struct {
	suite.Suite
	biletixSvc   service.ProviderService
	grpcConn     *grpc.ClientConn
	eventClient  pb.EventgRPCServiceClient
	amqqChannel  *amqp.Channel
	ampqService  mq.MQClient
	passoSvc     service.ProviderService
	schedulerSvc scheduler.SchedulerService
	httpClient   *resty.Client
	s            *gocron.Scheduler
	consumer     <-chan amqp.Delivery
	configs      *scraper.AppConfig
	v            *viper.Viper
}

func TestE2EScraperSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

type server struct {
	pb.UnimplementedEventgRPCServiceServer
}

// SendEvent is incoming event consumer on gRPC
func (s server) SendEvent(eventServer pb.EventgRPCService_SendEventServer) error {
	log.Infoln("SendEvent function  called by client is started")
	for {
		feature, err := eventServer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorln(err.Error())
			break
		}
		log.Println(feature)
		utils.FailOnError(err, "While creating event on DB error is received:")
	}
	return nil
}

func (ts *e2eScraperTestSuite) SetupSuite() {

	str := fmt.Sprintf("%s:%s", ts.configs.Grpc.ADDR, ts.configs.Grpc.PORT)

	ts.v = viper.New()
	_configs := scraper.NewConfig(ts.v)

	env := os.Getenv("ENV")
	appConfig, err := _configs.GetConfig(env)
	ts.configs = appConfig
	if err != nil {
		return
	}

	//gRPC connection
	_grpcConn, err := grpc.Dial(str, grpc.WithTransportCredentials(insecure.NewCredentials()))
	ts.grpcConn = _grpcConn
	utils.FailOnError(err, "did not connect")

	log.Infoln("gRPC connection is starting...")

	// RabbitMQ connection
	mqConn, err := amqp.Dial(ts.configs.RabbitMQ.AMPQ_URI)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	_amqqChannel, err := mqConn.Channel()
	ts.amqqChannel = _amqqChannel

	utils.FailOnError(err, "Failed to open a channel")

	log.Infoln("RabbitMQ connection is starting...")

	_eventClient := pb.NewEventgRPCServiceClient(_grpcConn)
	ts.eventClient = _eventClient

	//Gocron init
	_s := gocron.NewScheduler()
	ts.s = _s

	//Resty http client
	_httpClient := resty.New()
	ts.httpClient = _httpClient

	_ampqService := mq.NewAMQPService(_amqqChannel, log.New())
	_ampqService.ExchangeDeclare(ts.configs.RabbitMQ.EXCHANGE_NAME, "direct", ts.configs.RabbitMQ.QUEUE_NAME)
	ts.ampqService = _ampqService

	_biletixSvc := service.NewBiletixService(_eventClient, log.New(), _ampqService, _httpClient, ts.configs)
	ts.biletixSvc = _biletixSvc

	_passoSvc := service.NewPassoService(_eventClient, log.New(), _ampqService, _httpClient, ts.configs)
	ts.passoSvc = _passoSvc

	_schedulerSvc := scheduler.NewScheduler(_s)
	ts.schedulerSvc = _schedulerSvc

	_ampqService.QueueBind(ts.configs.RabbitMQ.QUEUE_NAME, ts.configs.RabbitMQ.ROUTING_KEY, ts.configs.RabbitMQ.EXCHANGE_NAME)
	_consumer := _ampqService.Consumer(ts.configs.RabbitMQ.QUEUE_NAME, ts.configs.RabbitMQ.CONSUMER_TAG, true)
	ts.consumer = _consumer

	// gRPC implementation
	s := grpc.NewServer()
	pb.RegisterEventgRPCServiceServer(s, &server{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", 50051))
	utils.FailOnError(err, "tcp listen failed.")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (ts *e2eScraperTestSuite) TestBiletixService_GetEvents() {
	var scheduleModels []model.ScheduleModel
	json.Unmarshal([]byte(ts.configs.SCHEDULE_ARRAY), &scheduleModels)

	scheduler.SetJob(ts.s,
		scheduleModels[0].TimeType,
		scheduleModels[0].TimeCount,
		scheduleModels[0].At,
		ts.biletixSvc.GetEvents,
		utils.GetQueryParameters(scheduleModels[0].Category, scheduleModels[0].City))

	ts.s.Start()
	go ts.ampqService.HandleConsumer(ts.consumer)

}
