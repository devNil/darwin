package main

import(
	"net/http"
	"os"
	"log"
	"darwin/web"
    "darwin/socket"
	ws "code.google.com/p/go.net/websocket"
)
const PORT = "8080"
func main(){
	var port string
	//This is just a test
	//This is just a test
	http.HandleFunc("/", web.IndexHandler)
	http.HandleFunc("/mobile", web.MobileIndexHandler)
	http.Handle("/wsb", ws.Handler(web.BrowserSocketHandler))
	http.HandleFunc("/register", web.RegisterMobileHandler)
	http.Handle("/wsm", ws.Handler(web.MobileSocketHandler))
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
