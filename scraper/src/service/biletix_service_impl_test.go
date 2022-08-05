package service

import (
	"encoding/json"
	"fmt"
	pb "github.com/cemayan/event-scraper-common/protos"
	"github.com/cemayan/event-scraper/scraper/src/config"
	"github.com/cemayan/event-scraper/scraper/src/model"
	"github.com/cemayan/event-scraper/scraper/src/mq"
	"github.com/cemayan/event-scraper/scraper/src/scheduler"
	"github.com/cemayan/event-scraper/scraper/src/utils"
	"github.com/go-resty/resty/v2"
	"github.com/jasonlvhit/gocron"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
	biletixSvc   ProviderService
	grpcConn     *grpc.ClientConn
	eventClient  pb.EventgRPCServiceClient
	amqqChannel  *amqp.Channel
	ampqService  mq.MQClient
	passoSvc     ProviderService
	schedulerSvc scheduler.SchedulerService
	httpClient   *resty.Client
	s            *gocron.Scheduler
	consumer     <-chan amqp.Delivery
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

var configs = config.GetConfig()

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

func (ts *e2eTestSuite) SetupSuite() {

	str := fmt.Sprintf("%s:%s", configs.GRPC_ADDR, configs.GRPC_ADDR_PORT)

	//gRPC connection
	_grpcConn, err := grpc.Dial(str, grpc.WithTransportCredentials(insecure.NewCredentials()))
	ts.grpcConn = _grpcConn
	utils.FailOnError(err, "did not connect")

	log.Infoln("gRPC connection is starting...")

	// RabbitMQ connection
	mqConn, err := amqp.Dial(configs.AMPQ_URI)
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
	_ampqService.ExchangeDeclare(configs.EXCHANGE_NAME, "direct", configs.QUEUE_NAME)
	ts.ampqService = _ampqService

	_biletixSvc := NewBiletixService(_eventClient, log.New(), _ampqService, _httpClient)
	ts.biletixSvc = _biletixSvc

	_passoSvc := NewPassoService(_eventClient, log.New(), _ampqService, _httpClient)
	ts.passoSvc = _passoSvc

	_schedulerSvc := scheduler.NewScheduler(_s)
	ts.schedulerSvc = _schedulerSvc

	_ampqService.QueueBind(configs.QUEUE_NAME, configs.ROUTING_KEY, configs.EXCHANGE_NAME)
	_consumer := _ampqService.Consumer(configs.QUEUE_NAME, configs.CONSUMER_TAG, true)
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

func (ts *e2eTestSuite) TestBiletixService_GetEvents() {
	var scheduleModels []model.ScheduleModel
	json.Unmarshal([]byte(configs.SCHEDULE_ARRAY), &scheduleModels)

	scheduler.SetJob(ts.s,
		scheduleModels[0].TimeType,
		scheduleModels[0].TimeCount,
		scheduleModels[0].At,
		ts.biletixSvc.GetEvents,
		utils.GetQueryParameters(scheduleModels[0].Category, scheduleModels[0].City))

	ts.s.Start()
	go ts.ampqService.HandleConsumer(ts.consumer)

}
