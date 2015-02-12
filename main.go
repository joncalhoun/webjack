package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	server := NewServer()
	http.Handle("/ws", server.GetHandler())
	http.Handle("/", http.FileServer(http.Dir("public")))
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var Connections []*websocket.Conn

func (s *Server) GetHandler() websocket.Handler {
	onConnect := func(ws *websocket.Conn) {
		Connections = append(Connections, ws)

		for {
			msg := map[string]string{
				"Message": fmt.Sprintf("%d connections!", len(Connections)),
			}
			err := websocket.JSON.Send(ws, msg)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(2 * time.Second)
		}
	}

	return websocket.Handler(onConnect)
}
