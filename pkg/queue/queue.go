package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	ConnectionString string
	Topic            string
	conn             *amqp.Connection
	channel          *amqp.Channel
}

func NewQueue(co string, topic string) *Queue {
	return &Queue{
		ConnectionString: co,
		Topic:            topic,
	}
}

func (q *Queue) Close() {
	if q.channel != nil {
		q.channel.Close()
		q.channel = nil
	}

	if q.conn != nil {
		q.conn.Close()
		q.conn = nil
	}
}

func (q *Queue) Send(ctx context.Context, res []byte) error {
	if q.channel == nil || q.conn == nil {
		q.Close()
		if err := q.Connect(); err != nil {
			return err
		}
	}

	d, err := q.channel.PublishWithDeferredConfirmWithContext(ctx, "", q.Topic, false, false, amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "image/jpeg",
		ContentEncoding: "",
		DeliveryMode:    amqp.Transient,
		Priority:        0,
		AppId:           "casoto",
		Body:            res,
	})
	if err != nil {
		return err
	}

	_, err = d.WaitContext(ctx)
	return err
}

func (q *Queue) Connect() error {
	config := amqp.Config{
		Vhost:      "/",
		Properties: amqp.NewConnectionProperties(),
	}
	config.Properties.SetClientConnectionName("catoso")

	var err error

	q.conn, err = amqp.DialConfig(q.ConnectionString, config)
	if err != nil {
		return err
	}

	q.channel, err = q.conn.Channel()
	if err != nil {
		q.Close()
		return err
	}

	_, err = q.channel.QueueDeclare(q.Topic, false, false, false, false, nil)
	if err != nil {
		q.Close()
		return err
	}

	if err := q.channel.Confirm(false); err != nil {
		q.Close()
		return err
	}

	return nil
}
