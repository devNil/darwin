package main

import (
	ws "code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"so/socket"
	"so/web"
)

func main() {
	http.HandleFunc("/", web.IndexHandler)

	http.HandleFunc("/js/", web.StaticHandler)

	server := socket.NewServer()
	server.Start()

	http.Handle("/ws", ws.Handler(server.Handler))

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}