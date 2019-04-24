package publisher

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
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

func (p *Publisher) connect() error {
	log.Info("[RabbitMQ] Reay to connect to server")

	conn, err := amqp.Dial(p.URI)
	p.Connection = conn

	go p.onConnectionClosed(conn)

	failOnError(err, "Failed to connect to RabbitMQ")
	return err
}

func (p *Publisher) onConnectionClosed(conn *amqp.Connection) {
	ch := conn.NotifyClose(make(chan *amqp.Error))
	go func() {
		defer close(ch)
		err := <-ch
		log.WithError(err).Error("Connection closed, should reconnect channel soon")

		p.connect()
	}()
}

func (p *Publisher) declareExchange() {
	ch, err := p.Connection.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()
	log.Info("[RabbitMQ] Forked a new channel")

	err = ch.ExchangeDeclare(
		p.ExchangeName,      // exchange name
		amqp.ExchangeFanout, // exchange type
		true,                // durable
		false,               // auto deleted
		false,               // internal
		false,               // no wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a exchange")
	p.Channel = ch
	log.Infof("[RabbitMQ] Declared a [%s] exchange -> %s", amqp.ExchangeFanout, p.ExchangeName)

	go p.onChannelClosed(ch)
}

func (p *Publisher) onChannelClosed(channel *amqp.Channel) {
	ch := channel.NotifyClose(make(chan *amqp.Error))
	go func() {
		defer close(ch)
		err := <-ch
		log.WithError(err).Error("channel closed, should reconnect channel soon")
		p.declareExchange()
	}()
}

func NewPublisher(uri string, exchange string) *Publisher {
	publisher := Publisher{URI: uri, ExchangeName: exchange}

	if err := publisher.connect(); err != nil {
		// reconect onece
		log.Info("[RabbitMQ] Reconnect in 5 second")
		time.Sleep(5 * time.Second)
		publisher.connect()
	}
	log.Info("[RabbitMQ] Connected successfully")

	publisher.declareExchange()
	return &publisher
}

func (p *Publisher) Publish(msgId, appId, topic string, body []byte) error {
	err := p.Channel.Publish(
		p.ExchangeName, // exchange
		"",             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "text/plian",
			DeliveryMode: amqp.Persistent,
			Priority:     0, // 0 to 9
			Timestamp:    time.Now(),
			MessageId:    msgId,
			AppId:        appId,
			Type:         topic,
			Body:         body,
		})

	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"MessageId": msgId,
			"AppId":     appId,
			"type":      topic,
			"Body":      string(body),
		}).Error("Failed to publish message")
	}
	return err
}

type AppPublisher struct {
	AppId     string
	Publisher *Publisher
}

func NewAppPublisher(uri, exchange, appId string) *AppPublisher {
	return &AppPublisher{
		Publisher: NewPublisher(uri, exchange),
		AppId:     appId,
	}
}

func (p *AppPublisher) Publish(topic string, payload interface{}) error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	msgId := uid.String()

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.Publisher.Publish(msgId, p.AppId, topic, body)
}
