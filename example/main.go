package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joncalhoun/webjack"
)

const longForm = "Jan 2, 2006 at 3:04:05pm (MST)"

func main() {
	server := webjack.NewServer()
	http.Handle("/ws", server.GetHandler())
	http.Handle("/", http.FileServer(http.Dir("public")))
	go ping(server)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func ping(server *webjack.Server) {
	c := time.Tick(5 * time.Second)
	for _ = range c {
		server.SendChannel(time.Now().Format(longForm), "test-channel")
	}
}
