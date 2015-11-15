package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os/user"

	"github.com/gorilla/mux"
)

// Load in the templates we need
var templates = template.Must(template.ParseFiles("templates/home.html"))

func main() {
	r := mux.NewRouter()

	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	// This handles static media
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static")),
	))

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		err := templates.ExecuteTemplate(w, "home.html", u.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":3000", r)
}
