package main

import (
	"fmt"
	pb "github.com/cemayan/event-scraper-common/protos"
	"github.com/cemayan/event-scraper/api/src/config"
	"github.com/cemayan/event-scraper/api/src/database"
	"github.com/cemayan/event-scraper/api/src/repo"
	"github.com/cemayan/event-scraper/api/src/utils"
	"github.com/sirupsen/logrus"
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
var configs = config.GetConfig()

func init() {

	//logrus init
	_log = logrus.New()
	_log.Out = os.Stdout

	// DB connection
	database.ConnectDB(configs)
	eventRepo = repo.NewEventRepo(database.GetDB(), _log)

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPC_PORT))
	utils.FailOnError(err, "tcp listen failed.")

	// gRPC implementation
	s := grpc.NewServer()
	pb.RegisterEventgRPCServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
