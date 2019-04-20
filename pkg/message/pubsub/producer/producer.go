package producer

import (
	"time"
	"github.com/streadway/amqp"
	log "github.com/sirupsen/logrus"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.WithError(err).Fatal(msg)
	}
}

type Producer struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	URI          string
	ExchangeName string
}

func NewProducer(uri string, exchange string) *Producer {
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()
	log.Info("Connect to RabbitMQ server successfully")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()
	log.Info("[RabbitMQ] forked a new channel")

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

	log.Infof("[RabbitMQ] Declared exchange[%s] -> %s", amqp.ExchangeFanout, exchange)

	p := Producer{
		URI:          uri,
		Connection:   conn,
		Channel:      ch,
		ExchangeName: exchange,
	}

	return &p
}


func (p *Producer) Publish(routingKey, msgAppId, msgType, MsgId string, msgBody []byte) error {
	return p.Channel.Publish(
			p.ExchangeName,  // exchange
			routingKey,      // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plian",
				DeliveryMode: amqp.Persistent,
				Priority:     0, // 0 to 9
				Timestamp:    time.Now(),
				AppId:        msgAppId,
				Type:         msgType,
				MessageId:    MsgId,
				Body:         msgBody,
			})
}