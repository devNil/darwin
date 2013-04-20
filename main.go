package main

import(
	"net/http"
	"os"
	"darwin/web"
	ws "code.google.com/p/go.net/websocket"
)

func main(){
	//This is just a test
	//This is just a test
	http.HandleFunc("/", web.IndexHandler)
	http.HandleFunc("/mobile", web.MobileIndexHandler)
	http.Handle("/wsb", ws.Handler(web.BrowserSocketHandler))
	http.HandleFunc("/register", web.RegisterMobileHandler)
	http.Handle("/wsm", ws.Handler(web.MobileSocketHandler))

	err:=http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil{
		panic(err)
	}
}
