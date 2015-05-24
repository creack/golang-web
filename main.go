package main

import (
	_ "expvar"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

// Page .
type Page struct {
	Title string
}

var (
	templates     = template.New("tpl")
	registeredTpl = map[string]struct{}{}
)

func init() {
	tpls, err := ioutil.ReadDir("tpl")
	if err != nil {
		log.Fatalf("Error reading template dir: %s", err)
	}
	for _, fi := range tpls {
		registeredTpl[fi.Name()] = struct{}{}
		template.Must(templates.ParseFiles(path.Join("tpl", fi.Name())))
	}
}

// RootHandler parses and serve content in ./tpl.
func RootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	if err := req.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %s", err), http.StatusInternalServerError)
	}
	path := mux.Vars(req)["path"]
	if path == "" || path == "/" {
		path = "index.tpl"
	}
	if !strings.HasSuffix(path, ".tpl") {
		path += ".tpl"
	}
	if _, ok := registeredTpl[path]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not found")
		return
	}
	if err := templates.ExecuteTemplate(w, path, Page{
		Title: "Home",
	}); err != nil {
		log.Printf("Error executing template: %s", err)
		http.Error(w, fmt.Sprintf("error parsing template: %s", err), http.StatusInternalServerError)
	}
}

func main() {
	println("Ready on :80")
	router := mux.NewRouter()
	router.Methods("GET").Path("/static/{path:.*}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, path.Join(".", "static", mux.Vars(req)["path"]))
	})
	router.Methods("GET").Path("/{path:.*}").HandlerFunc(RootHandler)
	log.Fatal(http.ListenAndServe(":80", router))
}
