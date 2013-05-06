package main

import (
	ws "code.google.com/p/go.net/websocket"
	"darwin/socket"
	"darwin/web"
	"log"
	"net/http"
	"os"
)

const PORT = "8080"

func main() {
	var port string
	http.HandleFunc("/", web.IndexHandler)
	http.HandleFunc("/mobile", web.MobileIndexHandler)
	http.HandleFunc("/register", web.RegisterMobileHandler)
	http.Handle("/ws", ws.Handler(socket.ConnectionHandler))
	http.Handle("/wsm", ws.Handler(socket.RemoteConnectionHandler))
    http.HandleFunc("/js/", web.JSSourceHandler)
    http.HandleFunc("/css/", web.CSSSourceHandler)
    http.HandleFunc("/img/", web.IMGSourceHandler)
	if os.Getenv("PORT") == "" {
		port = PORT
	} else {
		port = os.Getenv("PORT")
	}
	log.Println(port)
	socket.Run()
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
