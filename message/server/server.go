package main

import (
	"os"
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"

	_ "audi/pkg/logger"
	"audi/pkg/session/postgres"
	"audi/pkg/session/rabbitmq/pubsub/producer"
	pb "audi/message/proto"
	"audi/message/server/service"
	"audi/message/server/repository"
)

var (
	port     string
	exchange string
	PG_URL  string
	RABBITMQ_URL string
)

func init() {
	flag.StringVar(&port, "port", "3721", "The port of service listen")
	flag.StringVar(&exchange, "exchange", "pubsub", "The exchange name of rabbitmq")
	flag.Parse()

	initEnv()
}

func initEnv() {
	var ok bool
	PG_URL, ok = os.LookupEnv("PG_URL")
	if !ok {
		log.Fatal("Missing PG_URL setting")
	}

	RABBITMQ_URL, ok = os.LookupEnv("RABBITMQ_URL")

	if !ok {
		log.Fatal("Missing RABBITMQ_URL setting")
	}
}

func main() {
	db   := postgres.NewPostgres(PG_URL, true)
	msg  := producer.NewProducer(RABBITMQ_URL, exchange)
	repo := repository.NewRepository(db)
	server  := service.NewMessageServiceServer(repo, msg)

	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, server)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.WithError(err).Fatal("Failed to listen %s", addr)
	}

	log.Info("Initialized all components")
	log.Infof("Server starting with addr -> %s", addr)
	log.Fatal(s.Serve(l))
}