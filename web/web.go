//This package contains all web-specific contents like handler
//templates and structs
package web

import(
	"log"
	"os"
	"net/http"
	"html/template"
    "fmt"
)
const TEMPLATE = "template/"

var templates *template.Template

var jsPath string
var cssPath string
var imgPath string
func init(){
    var path string
	if os.Getenv("TEMPLATE") == "" {
		path = TEMPLATE
	} else {
		path = os.Getenv("TEMPLATE")
	}
	log.Println("Template-Directory: "+path)

    templates = template.Must(template.ParseGlob(fmt.Sprint(TEMPLATE,"*")))

    jsPath = path
    cssPath = path
    imgPath = path
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
    templates.ExecuteTemplate(w, "index", r.Host)
}

//Handler for the mobile-index-page
func MobileIndexHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "text/html")
    templates.ExecuteTemplate(w, "mobileindex", nil)
}
func JSSourceHandler(w http.ResponseWriter, r *http.Request) {
    log.Println(jsPath+r.URL.Path[1:])
    http.ServeFile(w, r, jsPath+r.URL.Path[1:])
}
func CSSSourceHandler(w http.ResponseWriter, r *http.Request) {
    log.Println(cssPath+r.URL.Path[1:])
    http.ServeFile(w, r, cssPath+r.URL.Path[1:])
}
func IMGSourceHandler(w http.ResponseWriter, r *http.Request) {
    log.Println(imgPath+r.URL.Path[1:])
    http.ServeFile(w, r, imgPath+r.URL.Path[1:])
}
func RegisterMobileHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm();

	id := r.FormValue("id")

	w.Header().Set("content-type", "text/html")
    
    templates.ExecuteTemplate(w, "controller", controllerPage{r.Host, id})
}
