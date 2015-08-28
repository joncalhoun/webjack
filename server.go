package webjack

import (
	"log"
	"net/url"

	"code.google.com/p/go.net/websocket"
)

type Server struct {
	channels    map[string]*Channel
	connections int
}

func NewServer() *Server {
	return &Server{
		make(map[string]*Channel),
		0,
	}
}

func (self *Server) NumConnections() int {
	return self.connections
}

func (self *Server) ConnectionsPerChannel() map[string]int {
	ret := make(map[string]int)
	for name, ch := range self.channels {
		ret[name] = ch.NumClients()
	}
	return ret
}

func (self *Server) AddClient(c *Client, ch string) {
	if self.channels[ch] == nil {
		self.channels[ch] = NewChannel()
	}
	self.channels[ch].AddClient(c)
	self.connections++
}

func (self *Server) RemoveClient(c *Client, ch string) {
	self.channels[ch].RemoveClient(c)
	self.connections--
}

func (self *Server) SendAll(msg interface{}) {
	log.Printf("Send to all channels with %+v\n", msg)
	for chName, _ := range self.channels {
		self.SendChannel(msg, chName)
	}
}

func (self *Server) SendChannel(msg interface{}, ch string) {
	if self.channels[ch] != nil {
		log.Printf("Sending %+v to channel %s\n", msg, ch)
		self.channels[ch].SendAll(msg)
	}
}

func (self *Server) Connections() int {
	return self.connections
}

func getChannelName(ws *websocket.Conn) string {
	u, err := url.Parse(ws.LocalAddr().String())
	if err == nil {
		v := u.Query()
		if val, ok := v["channel"]; ok {
			if len(val) > 0 {
				return val[0]
			}
		}
	}
	log.Println("Failed to find a channel. Using base channel of empty string.")
	// Base channel
	return ""
}

func (self *Server) GetHandler() websocket.Handler {
	onConnect := func(ws *websocket.Conn) {
		c := NewClient(ws)
		ch := getChannelName(ws)
		self.AddClient(c, ch)

		for {
			msg, exit := c.Listen()
			if exit {
				log.Printf("Removing client #%d\n", c.id)
				self.RemoveClient(c, ch)
				return
			} else if msg != nil {
				self.SendChannel(msg, ch)
			}
		}
	}

	return websocket.Handler(onConnect)
}
