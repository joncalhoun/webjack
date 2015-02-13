package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"log"
)

type Server struct {
	clients map[int]*Client
}

func NewServer() *Server {
	return &Server{
		make(map[int]*Client),
	}
}

func (s *Server) RemoveClient(c *Client) {
	delete(s.clients, c.id)
}

func (s *Server) AddClient(c *Client) {
	s.clients[c.id] = c
}

func (s *Server) SendAll(msg interface{}) {
	fmt.Printf("Send all with %+v\n", msg)
	for _, c := range s.clients {
		err := c.Send(msg)
		if err != nil {
			fmt.Printf("Error sending to client #%d: %s\n", c.id, err)
		}
	}
}

func (s *Server) GetHandler() websocket.Handler {
	onConnect := func(ws *websocket.Conn) {
		c := NewClient(ws)
		s.AddClient(c)

		for {
			select {
			default:
				var msg interface{}
				err := c.Receive(&msg)
				if err == io.EOF {
					log.Printf("Removing client #%d\n", c.id)
					s.RemoveClient(c)
					return
				} else if err != nil {
					log.Println(err)
				} else {
					log.Printf("Received from client #%d: %+v\n", c.id, msg)
					s.SendAll(msg)
				}
			}
		}
	}

	return websocket.Handler(onConnect)
}
