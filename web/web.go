//This package contains all web-specific contents like handler
//templates and structs
package web

import(
	"log"
	"os"
	"net/http"
	"html/template"
)
const TEMPLATE = "web/"
var index *template.Template
var mobileIndex *template.Template
var controller *template.Template

func init(){
	var path string
	if os.Getenv("TEMPLATE") == "" {
		path = TEMPLATE
	} else {
		path = os.Getenv("TEMPLATE")
	}
	log.Println("Template-Directory: "+path)
	index = template.Must(template.ParseFiles(path+"index.html"))
	mobileIndex = template.Must(template.ParseFiles(path+"mobile.html"))
	controller = template.Must(template.ParseFiles(path+"controller.html"))
}

//Struct for the index template
type indexPage struct{
	Url string
}

type controllerPage struct{
	Url string
	Id string
}

//Handler for the index-Page
func IndexHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "text/html")
	index.Execute(w, indexPage{r.Host})
}

//Handler for the mobile-index-page
func MobileIndexHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "text/html")
	mobileIndex.Execute(w, nil)
}

func RegisterMobileHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm();

	id := r.FormValue("id")

	w.Header().Set("content-type", "text/html")

	controller.Execute(w, controllerPage{r.Host, id})
}
