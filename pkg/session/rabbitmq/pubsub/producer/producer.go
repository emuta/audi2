package producer

import (
	"time"
	"strconv"

	"github.com/streadway/amqp"
	log "github.com/sirupsen/logrus"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.WithError(err).Fatal(msg)
	}
}

type Message struct {
	AppId string
	Type  string
	Id    int64
	Body  []byte
}

func NewMessage(appId string, tp string, id int64, body []byte) *Message {
	return &Message{AppId: appId, Type: tp, Id: id, Body: body}
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

	p := Producer{
		URI:          uri,
		Connection:   conn,
		Channel:      ch,
		ExchangeName: exchange,
	}

	return &p
}

/*
func convertDataToBytes(data interface{}) ([]byte, error) {
	if v, ok := data.(string); ok {
		return []byte(v), nil
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).WithField("data", data).Error("Failed to Marshal message")
		return nil, err
	}
	return b, nil
}
*/

func (p *Producer) Publish(msg *Message) error {

	return p.Channel.Publish(
			p.ExchangeName,  // exchange
			"",              // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plian",
				DeliveryMode: amqp.Persistent,
				Priority:     0, // 0 to 9
				Timestamp:    time.Now(),
				MessageId:    strconv.FormatInt(msg.Id, 10),
				AppId:        msg.AppId,
				Type:         msg.Type,
				Body:         msg.Body,
			})
}

/*
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
*/