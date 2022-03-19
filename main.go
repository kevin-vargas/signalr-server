package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/philippseith/signalr"
)

func main() {
	godotenv.Load(".env")
	address := "0.0.0.0:8081"

	hub := &AppHub{}

	server, _ := signalr.NewServer(context.TODO(),
		signalr.SimpleHubFactory(hub),
		signalr.KeepAliveInterval(2*time.Second),
	)

	router := http.NewServeMux()

	server.MapHTTP(signalr.WithHTTPServeMux(router), "/chat")
	fmt.Printf("Listening for websocket connections on http://%s\n", address)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:6006", "http://fast.ar:6006"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Requested-With", "X-Signalr-User-Agent"},
	})
	handler := c.Handler(router)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
