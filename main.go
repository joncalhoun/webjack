package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	server := NewServer()
	http.Handle("/ws", server.GetHandler())
	http.Handle("/", http.FileServer(http.Dir("public")))

	go PingSockets(server)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type PingMessage struct {
	Message string `json:"message"`
}

func PingSockets(s *Server) {
	for {
		time.Sleep(5 * time.Second)

		msg := PingMessage{
			Message: fmt.Sprintf("%d connections!", s.Connections()),
		}
		log.Printf("Sending: %+v\n", msg)

		s.SendAll(msg)
	}
}
