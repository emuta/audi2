package publisher

import (
	"time"
	"encoding/json"

	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	log "github.com/sirupsen/logrus"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.WithError(err).Fatal(msg)
	}
}

type Publisher struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	URI          string
	ExchangeName string
}

func NewPublisher(uri string, exchange string) *Publisher {
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()
	log.Info("Connected to RabbitMQ server")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()
	log.Info("[RabbitMQ] Forked a new channel")

	err = ch.ExchangeDeclare(
		exchange,            // exchange name
		amqp.ExchangeFanout, // exchange type
		true,                // durable
		false,               // auto deleted
		false,               // internal
		false,               // no wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	// log.Infof("[RabbitMQ] Declared exchange -> %s", exchange)
	log.Infof("[RabbitMQ] Declared a [%s] exchange -> %s", amqp.ExchangeFanout, exchange)

	p := Publisher{
		URI:          uri,
		Connection:   conn,
		Channel:      ch,
		ExchangeName: exchange,
	}

	return &p
}

func (p *Publisher) Publish(appId, tp string, payload interface{}) error {

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	uidStr := uid.String()

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.Channel.Publish(
			p.ExchangeName,  // exchange
			"",              // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plian",
				DeliveryMode: amqp.Persistent,
				Priority:     0, // 0 to 9
				Timestamp:    time.Now(),
				MessageId:    uidStr,
				AppId:        appId,
				Type:         tp,
				Body:         body,
			})

	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"MessageId": uidStr,
			"AppId": appId,
			"type":  tp,
			"Body":  string(body),
		}).Error("Failed to publish event")
	}
	return err
}

type AppPublisher struct {
	AppId      string
	Publisher *Publisher
}

func NewAppPublisher(uri, exchange, appId string) *AppPublisher{
	return &AppPublisher{
		Publisher: NewPublisher(uri, exchange),
		AppId: appId,
	}
}

func (p *AppPublisher) Publish(tp string, data interface{}) error {
	return p.Publisher.Publish(p.AppId, tp, data)
}


