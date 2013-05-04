//This package contains all web-specific contents like handler
//templates and structs
package web

import(
    "io/ioutil"
	"log"
	"os"
	"net/http"
	"html/template"
)
const TEMPLATE = "web/"
var index *template.Template
var mobileIndex *template.Template
var controller *template.Template
var gamejs []byte

func init(){
	var path string
	if os.Getenv("TEMPLATE") == "" {
		path = TEMPLATE
	} else {
		path = os.Getenv("TEMPLATE")
	}
	log.Println("Template-Directory: "+path)
	index = template.Must(template.ParseFiles(path+"game.html"))
	mobileIndex = template.Must(template.ParseFiles(path+"mobile.html"))
	controller = template.Must(template.ParseFiles(path+"controller.html"))
    gamejs,_ = ioutil.ReadFile(path+"game.js")
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
	index.Execute(w, r.Host)
}

//Handler for the mobile-index-page
func MobileIndexHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "text/html")
	mobileIndex.Execute(w, nil)
}
func GameHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("content-type", "application/javascript")
    w.Write(gamejs)
}

func RegisterMobileHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm();

	id := r.FormValue("id")

	w.Header().Set("content-type", "text/html")

	controller.Execute(w, controllerPage{r.Host, id})
}
