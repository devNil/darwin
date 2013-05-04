package main

import(
	"net/http"
	"os"
	"log"
	"darwin/web"
	ws "code.google.com/p/go.net/websocket"
	"darwin/socket"
)
const PORT = "8080"
func main(){
	var port string
	//This is just a test
	http.HandleFunc("/", web.IndexHandler)
	http.HandleFunc("/game.js", web.GameHandler)
	http.HandleFunc("/mobile", web.MobileIndexHandler)
	http.HandleFunc("/register", web.RegisterMobileHandler)
	http.Handle("/ws", ws.Handler(socket.ConnectionHandler))
	if os.Getenv("PORT") == "" {
		port = PORT
	} else {
		port = os.Getenv("PORT")
	}
	log.Println( port )
    go socket.Run()
	err:=http.ListenAndServe(":"+port, nil)
	if err != nil{
		panic(err)
	}
}
