package main

import(
	"net/http"
	"os"
)

func main(){
	http.HandleFunc("/", func(r http.ResponseWriter, req *http.Request){
		r.Write([]byte("Hello World"))
		return
	})

	err:=http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil{
		panic(err)
	}
}
