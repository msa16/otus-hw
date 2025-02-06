package client

import "context"

// интерфейс клиента кафки чтобы их можно было менять.
type Broker interface {
	Connect(ctx context.Context) error
	Disconnect() error
	Publish(topic string, message []byte) error
	Subscribe(topic string) (<-chan []byte, error)
}
