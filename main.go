package main

import(
	"net/http"
	"os"
)

func main(){
	//This is just a test
	//This is just a test
	http.HandleFunc("/", func(r http.ResponseWriter, req *http.Request){
		r.Write([]byte("Hello World"))
		return
	})

	err:=http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil{
		panic(err)
	}
}
