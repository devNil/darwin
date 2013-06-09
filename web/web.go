//This package mainly serves the templates
package web

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*"))
}

//Handler for the normal user, not mobile
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	name, err := os.Hostname()

	if err != nil {
		log.Println(err)
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Println(err)
		return
	}

	values := map[string]interface{}{"host": fmt.Sprint(addrs[0], ":8080")}
	values["wshost"] = template.HTML(fmt.Sprint("ws://", addrs[0], ":8080/ws"))

	w.Header().Add("content-type", "text/html")
	templates.ExecuteTemplate(w, "desktop", values)
	return
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	http.ServeFile(w, r, fmt.Sprint("static", r.URL.Path))
}
