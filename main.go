package main

import (
	ws "code.google.com/p/go.net/websocket"
	"darwin/socket"
	"darwin/web"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", web.IndexHandler)

	http.HandleFunc("/js/", web.StaticHandler)

	server := socket.NewServer()
	server.Start()

	http.Handle("/ws", ws.Handler(server.Handler))

	log.Println("Server started on :9080")
	err := http.ListenAndServe(":9080", nil)
	if err != nil {
		panic(err)
        log.Println("test")
	}
}
