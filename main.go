package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	server := NewServer()
	http.Handle("/ws", server.GetHandler())
	http.Handle("/", http.FileServer(http.Dir("public")))

	go PingSockets()
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var Connections []*websocket.Conn

func PingSockets() {
	for {
		time.Sleep(5 * time.Second)

		msg := map[string]string{
			"Message": fmt.Sprintf("%d connections!", len(Connections)),
		}
		for _, c := range Connections {
			err := websocket.JSON.Send(c, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func removeConnection(ws *websocket.Conn) {
	newConn := []*websocket.Conn{}
	for _, c := range Connections {
		if c != ws {
			newConn = append(newConn, c)
		}
	}
	Connections = newConn
}

func (s *Server) GetHandler() websocket.Handler {
	onConnect := func(ws *websocket.Conn) {
		Connections = append(Connections, ws)

		for {
			select {
			default:
				var msg struct {
					Message string `json:"message"`
				}
				err := websocket.JSON.Receive(ws, &msg)
				if err == io.EOF {
					removeConnection(ws)
					return
				} else if err != nil {
					log.Println(err)
				} else {
					log.Printf("Received: %+v\n", msg)
				}
			}
		}
	}

	return websocket.Handler(onConnect)
}
