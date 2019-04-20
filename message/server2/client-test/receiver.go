package main

import (
	"time"
	"net"
	// "encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const addr = "amqp://guest:guest@rabbitmq:5672/"
const exName = "pubsub"

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.WithError(err).Fatal(msg)
	}
}

func main() {
	conn, err := amqp.DialConfig(addr, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout("tcp", addr, 2 * time.Second)
		},
	})
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exName, // exchange name
		amqp.ExchangeFanout, // exchange type
		true,  // durable
		false, // auto deleted
		false, // internal
		false, // no wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a exchange")


	q, err := ch.QueueDeclare(
		"pubsub.receiver",    // queue name
		true, // durable
		false, //delete when unused
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"pusbsub.all",     // routing key
		exName, // exchange name
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",    // consumer
		true,  // audo ack
		false, //exclusive
		false, // no local
		false, // no wait
		nil,   // arguments
	)

	forever := make(chan bool)
	go func() {
		defer close(forever)

		for msg := range msgs {
			log.WithFields(log.Fields{
				"RoutingKey": msg.RoutingKey,
				"AppId":      msg.AppId,
				"Type":       msg.Type,
				"MessageId":   msg.MessageId,
			}).Infof("Received a message: %s", msg.Body)
			/*
			j, err := json.Marshal(msg)
			failOnError(err, "Fail to marshal message")
			log.Info(string(j))
			*/
		}
	}()
	<- forever
}