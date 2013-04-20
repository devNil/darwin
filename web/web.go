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
var mobileIndex *template.Template

func init(){
	log.Println("Template-Directory: "+os.Getenv("TEMPLATE"))
	path := os.Getenv("TEMPLATE")
	index = template.Must(template.ParseFiles(path+"index.html"))
	mobileIndex = template.Must(template.ParseFiles(path+"mobile.html"))
}

//Struct for the index template
type indexPage struct{
	Code int64
}

//Handler for the index-Page
func IndexHandler(w http.ResponseWriter, r *http.Request){
	unix := time.Now().Unix()
	w.Header().Set("content-type", "text/html")
	index.Execute(w, indexPage{unix})
}

//Handler for the mobile-index-page
func MobileIndexHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "text/html")
	mobileIndex.Execute(w, nil)
}
