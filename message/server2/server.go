package server

import (
	"os"
	"net"

	"google.golang.org/grpc"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"

	pb "audi/protobuf/message"
	"audi/message/pubsub/producer"
	"audi/message/server/service"
)
const (
	DefaultRabbitMQURI      = "amqp://guest:guest@rabbitmq:5672/"
	DefaultRabbitMQExchange = "pubsub"
	DefaultPostgreSQLURI    = "postgresql://postgres:postgres@postgresql:5432/postgres?sslmode=disable"
)

var (
	db *gorm.DB

	RabbitMQ_URI      string
	RabbitMQ_Exchange string
)

func init() {

	initLogger()

	initPostgres()
	log.Info("Connet Postgresql server ok")

	initRabbitMQ()
	log.WithField("exchange", RabbitMQ_Exchange).Info("Connet RabbitMQ   server ok")
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func initRabbitMQ() {
	uri, ok := os.LookupEnv("RABBITMQ_URL")
	if !ok {
		uri = DefaultRabbitMQURI
		log.WithField("RABBITMQ_URL", uri).Warn("RABBITMQ_URL not found in enviroment, use default")
	}
	RabbitMQ_URI = uri
	
	initRabbitMQExchange()
}

func initRabbitMQExchange() {
	name, ok := os.LookupEnv("RABBITMQ_EXCHANGE")
	if !ok {
		name = DefaultRabbitMQExchange
		log.WithField("RABBITMQ_EXCHANGE", name).Warn("RABBITMQ_EXCHANGE not found in enviroment, use default")
	}
	RabbitMQ_Exchange = name
}

func initPostgres() {
	var err error
	uri, ok := os.LookupEnv("PG_URL")
	if !ok {
		uri := DefaultPostgreSQLURI
		log.WithField("PG_URL", uri).Warn("PG_URL not found in enviroment, use default")
	}

	db, err = gorm.Open("postgres", uri)
    if err != nil {
        log.WithError(err).Fatal("Unable to connect database")
    }
    // db.SetLogger(log.StandardLogger())
    db.LogMode(true)
}

func Serve(addr string) {
	messageServiceServer := service.NewMessageServiceServer(db, producer.NewProducer(RabbitMQ_URI, RabbitMQ_Exchange))
	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, messageServiceServer)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.WithError(err).Fatalf("Failed to listen address: %s", addr)
	}

	log.WithField("listen", addr).Info("Server starting")

	if err := s.Serve(l); err != nil {
		log.WithError(err).Fatalf("Failed to boot grpc server")
	}
}