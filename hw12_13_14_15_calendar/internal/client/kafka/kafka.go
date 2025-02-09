package kafka

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/watermill"                              //nolint:depguard
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"           //nolint:depguard
	"github.com/ThreeDotsLabs/watermill/message"                      //nolint:depguard
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"    //nolint:depguard
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"        //nolint:depguard
	"github.com/google/uuid"                                          //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"    //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/client" //nolint:depguard
)

type Client struct {
	brokers    []string
	marshaler  kafka.DefaultMarshaler
	logger     watermill.LoggerAdapter
	publisher  message.Publisher
	subscriber message.Subscriber
	handlers   map[string]client.SubscriberHandlerFunc
	appLogger  app.Logger
}

var ErrPublisherNotReady = errors.New("publisher is not ready")

func New(brokers []string, appLogger app.Logger) *Client {
	return &Client{
		brokers: brokers, marshaler: kafka.DefaultMarshaler{}, logger: watermill.NewStdLogger(true, false),
		handlers: make(map[string]client.SubscriberHandlerFunc, 0), appLogger: appLogger,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	for {
		// ТЗ: при запуске процесс подключается к Kafka, если Kafka недоступна - ждёт
		// но надо останавливаться по Ctrl+C
		// docker run -rm -p 9092:9092 apache/kafka-native:3.9.0
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
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

	for topic, handler := range c.handlers {
		router.AddNoPublisherHandler(
			"handler_"+topic, // handler name, must be unique
			topic,            // topic from which messages should be consumed
			c.subscriber,
			func(msg *message.Message) error {
				raw := []byte(msg.Payload)
				err := handler(&raw)
				if err != nil {
					// When a handler returns an error, the default behavior is to send a Nack (negative-acknowledgement).
					// The message will be processed again.
					//
					// You can change the default behaviour by using middlewares, like Retry or PoisonQueue.
					// You can also implement your own middleware.
					c.appLogger.Error("handler" + err.Error())
					return err
				}
				return nil
			},
		)
	}

	if err := router.Run(ctx); err != nil {
		// это не ошибка коннекта, работать нельзя
		return err
	}
	return nil
}

func (c *Client) Disconnect() error {
	return nil
}

func (c *Client) Publish(topic string, payload []byte) error {
	if c.publisher == nil {
		return ErrPublisherNotReady
	}
	return c.publisher.Publish(topic, message.NewMessage(
		uuid.New().String(),
		payload,
	))
}

func (c *Client) Subscribe(topic string, handler client.SubscriberHandlerFunc) {
	c.handlers[topic] = handler
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
