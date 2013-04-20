package main

import(
	"net/http"
	"os"
	"darwin/web"
)

func main(){
	//This is just a test
	//This is just a test
	http.HandleFunc("/", web.IndexHandler)
	http.HandleFunc("/mobile", web.MobileIndexHandler)

	err:=http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil{
		panic(err)
	}
}
