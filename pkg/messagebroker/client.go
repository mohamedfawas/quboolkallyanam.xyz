package messaging

type MessageHandler func([]byte) error

type Client interface {
	Publish(topic string, data interface{}) error
	Subscribe(topic string, handler MessageHandler) error
	Close() error
}
