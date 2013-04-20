//This package contains all web-specific contents like handler
//templates and structs
package web

import(
	"log"
	"os"
	"net/http"
	"html/template"
	"time"
)

var index *template.Template

func init(){
	log.Println("Template-Directory: "+os.Getenv("TEMPLATE"))
	path := os.Getenv("TEMPLATE")
	index = template.Must(template.ParseFiles(path+"index.html"))
}

//Struct for the index template
type indexPage struct{
	Code int64
}

func IndexHandler(w http.ResponseWriter, r *http.Request){
	unix := time.Now().Unix()
	w.Header().Set("content-type", "text/html")
	index.Execute(w, indexPage{unix})
}

