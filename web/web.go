//This package mainly serves the templates
package web

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
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
    var addr string

	addrs, err := net.LookupHost(name)
	if err != nil {
        addr = "localhost"
		log.Println(err)
	} else {
        addr = addrs[0]
    }

	values := map[string]interface{}{"host": fmt.Sprint(addr, ":9080")}
	values["wshost"] = template.HTML(fmt.Sprint("ws://", addr, ":9080/ws"))

	w.Header().Add("content-type", "text/html")

	if strings.Contains(r.UserAgent(), "iPhone") {
		templates.ExecuteTemplate(w, "mobile", values)
		return
	}

	if strings.Contains(r.UserAgent(), "iPad") {
		templates.ExecuteTemplate(w, "mobile", values)
		return
	}

	templates.ExecuteTemplate(w, "desktop", values)
	return
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	http.ServeFile(w, r, fmt.Sprint("static", r.URL.Path))
}
