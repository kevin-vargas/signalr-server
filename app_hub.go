package main

import (
	"encoding/json"
	"fmt"

	"github.com/kevin-vargas/signalr-server/constants"
	"github.com/kevin-vargas/signalr-server/pubsub"
	"github.com/philippseith/signalr"
)

type AppHub struct {
	signalr.Hub
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
	return nil
}

// dont need
func (h *AppHub) ReceiveFromTopic(topic string, payload []byte) {
	h.Clients().Group(topic).Send(constants.RECEIVE_METHOD, payload)
}
