package main

import (
	"encoding/json"
	"flag"
	"fmt"

	_ "audi/pkg/logger"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	brokerURL    string
	exchangeName string

	queueName  string
	routingKey string
)

func failOnError(err error, msg string) {
	if err != nil {
		log.WithError(err).Fatal(msg)
	}
}

func init() {
	flag.StringVar(&brokerURL, "broker", "amqp://guest:guest@rabbitmq:5672/", "The broker url for RabbitMQ connect")
	flag.StringVar(&exchangeName, "exchange", "pubsub", "The exchange name of RabbitMQ used")
	flag.StringVar(&queueName, "queue", "queue", "The exchange name of RabbitMQ used")
	flag.StringVar(&routingKey, "routing", "routing", "The exchange name of RabbitMQ used")
	flag.Parse()

	queueName = fmt.Sprintf("%s.monitor.%s", exchangeName, queueName)
	routingKey = fmt.Sprintf("%s.monitor.%s", exchangeName, routingKey)
}

func watchConnection(conn *amqp.Connection) {
	ch := conn.NotifyClose(make(chan *amqp.Error))
	go func() {
		defer close(ch)
		err := <-ch
		log.WithError(err).Info("Connection closed, should reconnect channel soon")
	}()
}

func watchChannel(channel *amqp.Channel) {
	ch := channel.NotifyClose(make(chan *amqp.Error))
	go func() {
		defer close(ch)
		err := <-ch
		log.WithError(err).Info("channel closed, should reconnect channel soon")
	}()
}

func main() {
	conn, err := amqp.Dial(brokerURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	go watchConnection(conn)

	log.Info("Connected to RabbitMQ server")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	go watchChannel(ch)

	log.Info("[RabbitMQ] Forked a new channel")

	err = ch.ExchangeDeclare(
		exchangeName,        // exchange name
		amqp.ExchangeFanout, // exchange type
		true,                // durable
		false,               // auto deleted
		false,               // internal
		false,               // no wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	log.Infof("[RabbitMQ] Declared a [%s] exchange -> %s", amqp.ExchangeFanout, exchangeName)

	q, err := ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     //delete when unused
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Infof("[RabbitMQ] Declared a queue -> %s", queueName)

	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange name
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	log.Infof("[RabbitMQ] Bind queue with routingKey -> %s", routingKey)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // audo ack
		false,  //exclusive
		false,  // no local
		false,  // no wait
		nil,    // arguments
	)

	forever := make(chan bool)
	go func() {
		defer close(forever)

		for msg := range msgs {
			/*
				log.WithFields(log.Fields{
					"RoutingKey": msg.RoutingKey,
					"AppId":      msg.AppId,
					"Type":       msg.Type,
				}).Infof("Received a message: %s", msg.Body)
			*/
			data, err := json.Marshal(msg)
			if err != nil {
				log.Error(err)
			}
			log.Infof("%s", data)
		}
	}()
	<-forever
}
