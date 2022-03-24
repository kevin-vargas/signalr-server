package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kevin-vargas/signalr-server/constants"
	"github.com/kevin-vargas/signalr-server/pubsub"
	"github.com/kevin-vargas/signalr-server/topics"
	"github.com/philippseith/signalr"
)

type AppHub struct {
	signalr.Hub
}

var once sync.Once
var cancel = make(chan bool, 0)

func InitNotifyTopics(h *AppHub) {
	once.Do(func() {
		go h.NotifyTopicsUpdates(cancel)
	})
}

func (h *AppHub) Subscribe(topic string) error {
	id := h.ConnectionID()
	h.Groups().AddToGroup(topic, id)
	//TODO: IDEMPOTENTE
	err := pubsub.Get().Subscribe(topic, &h.Hub)
	if err != nil {
		fmt.Println("ON SUBSCRIBE", err)
		return err
	}
	fmt.Println("ALGUIEN SE SUBSCRIBIO")
	return nil
}

func (h *AppHub) Publish(topic string, payload json.RawMessage) error {
	err := pubsub.Get().Publish(topic, payload)
	if err != nil {
		fmt.Println("ON PUBLISH", err)
		return err
	}
	fmt.Println("SE PUBLICO ALGO")
	topics.Get().UpdateTopic(topic)
	return nil
}
func (h *AppHub) OnConnected(connectionID string) {
	InitNotifyTopics(h)
	fmt.Println(connectionID, " connected")
}

func (h *AppHub) NotifyTopicsUpdates(cancel <-chan bool) {
	for {
		select {
		case <-time.After(1 * time.Second):
			topics.Get().GetValids()
			h.Hub.Clients().All().Send("topics", topics.Get().GetValids())
		case <-cancel:
			return
		}
	}
}

// dont need
func (h *AppHub) ReceiveFromTopic(topic string, payload []byte) {
	h.Clients().Group(topic).Send(constants.RECEIVE_METHOD, payload)
}

func (h *AppHub) ValidTopics() {
	id := h.ConnectionID()
	topics := topics.Get().GetValids()
	h.Clients().Client(id).Send(constants.TOPICS_METHOD, topics)
}
