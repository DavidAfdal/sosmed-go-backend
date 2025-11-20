package rabbitmq

import (
	"github.com/davidafdal/post-app/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	Consume(queue string) (<-chan amqp.Delivery, error)
	Publish(exchange, routingKey, eventType string, body []byte) error
}

type Client struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewClient(cfg *config.RabbitConfig) (*Client, error) {
	conn, err := amqp.Dial(cfg.Url)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(
		queue, "", true, false, false, false, nil,
	)
}

func (r *Client) Publish(exchange, routingKey, eventType string, body []byte) error {
	return r.Channel.Publish(
		exchange, routingKey, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Type:         eventType,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (r *Client) Close() {
	r.Channel.Close()
	r.Conn.Close()
}
