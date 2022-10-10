package main

import (
	"fmt"
	"github.com/cemayan/event-scraper/config/api"
	"github.com/cemayan/event-scraper/internal/api/database"
	"github.com/cemayan/event-scraper/internal/api/repo"
	"github.com/cemayan/event-scraper/internal/api/utils"
	pb "github.com/cemayan/event-scraper/protos"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

type server struct {
	pb.UnimplementedEventgRPCServiceServer
}

var eventRepo repo.EventRepository
var _log *logrus.Logger
var configs *api.AppConfig
var v *viper.Viper
var dbHandler database.DBHandler

func init() {

	//logrus init
	_log = logrus.New()
	_log.Out = os.Stdout

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

	_log.Infoln("RabbitMQ connection is starting...")
}

// SendEvent is incoming event consumer on gRPC
func (s server) SendEvent(eventServer pb.EventgRPCService_SendEventServer) error {
	_log.Infoln("SendEvent function  called by client is started")
	for {
		feature, err := eventServer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			_log.Errorln(err.Error())
			break
		}
		_, err = eventRepo.Create(feature)
		utils.FailOnError(err, "While creating event on DB error is received:")
	}
	return nil
}

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.Grpc.PORT))
	utils.FailOnError(err, "tcp listen failed.")

	// gRPC implementation
	s := grpc.NewServer()
	pb.RegisterEventgRPCServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
