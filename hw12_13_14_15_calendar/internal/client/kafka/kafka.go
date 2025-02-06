package kafka

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"                           //nolint:depguard
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"        //nolint:depguard
	"github.com/ThreeDotsLabs/watermill/message"                   //nolint:depguard
	"github.com/ThreeDotsLabs/watermill/message/router/middleware" //nolint:depguard
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"     //nolint:depguard
)

type Client struct {
	brokers    []string
	marshaler  kafka.DefaultMarshaler
	logger     watermill.LoggerAdapter
	publisher  message.Publisher
	subscriber message.Subscriber
}

func New(brokers []string) Client {
	return Client{brokers: brokers, marshaler: kafka.DefaultMarshaler{}, logger: watermill.NewStdLogger(true, false)}
}

func (c *Client) Connect(ctx context.Context) error {
	for {
		// пробуем подключаться пока не получиться
		var err error
		c.publisher, err = c.createPublisher()
		if err != nil {
			c.logger.Error("createPublisher", err, nil)
			continue
		}
		// Subscriber is created with consumer group handler_1
		c.subscriber, err = c.createSubscriber("handler_1")
		if err == nil {
			break
		}
		c.logger.Error("createSubscriber", err, nil)
		err = c.publisher.Close()
		if err != nil {
			c.logger.Error("closePublisher", err, nil)
		}
	}

	router, err := message.NewRouter(message.RouterConfig{}, c.logger)
	if err != nil {
		// это не ошибка коннекта, работать нельзя
		return err
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(middleware.Recoverer)

	if err := router.Run(ctx); err != nil {
		// это не ошибка коннекта, работать нельзя
		return err
	}
	return nil
}

func (c *Client) Disconnect() error {
	return nil
}

func (c *Client) Publish(_ string, _ []byte) error {
	return nil
}

func (c *Client) Subscribe(_ string) (<-chan []byte, error) {
	return nil, nil
}

func (c *Client) createPublisher() (message.Publisher, error) {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   c.brokers,
			Marshaler: c.marshaler,
		},
		c.logger,
	)
	if err != nil {
		return nil, err
	}
	return kafkaPublisher, nil
}

func (c *Client) createSubscriber(consumerGroup string) (message.Subscriber, error) {
	kafkaSubscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       c.brokers,
			Unmarshaler:   c.marshaler,
			ConsumerGroup: consumerGroup, // every handler will use a separate consumer group
		},
		c.logger,
	)
	if err != nil {
		return nil, err
	}
	return kafkaSubscriber, nil
}
