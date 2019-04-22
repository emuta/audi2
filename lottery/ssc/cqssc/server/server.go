package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	_ "audi/pkg/logger"
	"audi/pkg/session/postgres"
	"audi/pkg/session/rabbitmq/pubsub/publisher"
	pb "audi/lottery/ssc/cqssc/proto"
	"audi/lottery/ssc/cqssc/server/repository"
	"audi/lottery/ssc/cqssc/server/service"
	
)

const ServiceName = "lottery.ssc.cqssc"

var (
	port     string
	exchange string
)

func init() {
	flag.StringVar(&port, "port", "3721", "The port of service listen")
	flag.StringVar(&exchange, "exchange", "pubsub", "The exchange name of RabbitMQ used")
	flag.Parse()
}

func main() {
	db := postgres.NewPostgres(os.Getenv("PG_URL"), true)
	puber := publisher.NewAppPublisher(os.Getenv("RABBITMQ_URL"), exchange, ServiceName)
	repo := repository.NewRepository(db)
	server := service.NewCqsscServiceServer(repo, puber)

	s := grpc.NewServer()
	pb.RegisterCqsscServiceServer(s, server)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.WithError(err).Fatal("Failed to listen %s", addr)
	}

	log.Info("Initialized all components")
	log.Infof("Server starting with addr -> %s", addr)
	log.Fatal(s.Serve(l))
}
