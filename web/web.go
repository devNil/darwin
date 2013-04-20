//This package contains all web-specific contents like handler
//templates and structs
package web

import(
	"log"
	"os"
	"net/http"
	"html/template"
)

var index *template.Template

func init(){
	log.Println(os.Getenv("TEMPLATE"))
	path := os.Getenv("TEMPLATE")
	index = template.Must(template.ParseFiles(path+"index.html"))
}

func IndexHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "text/html")
	index.Execute(w, nil)
	return
}

