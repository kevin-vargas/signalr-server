package pubsub

import (
	"encoding/json"
	"sync"

	"github.com/philippseith/signalr"
)

type Result struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

type Client interface {
	Publish(topic string, payload []byte) error
	SubscribeDefault(topic string) error
	Subscribe(topic string, hub *signalr.Hub) error
}

func GetMQTTClient() Client {
	return NewMQTTClient()
}

var once sync.Once
var instance Client

func Get() Client {
	once.Do(func() {
		client := GetMQTTClient()
		instance = client
	})
	return instance
}
