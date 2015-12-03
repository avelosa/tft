package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os/user"
	"os"
	"io"

	"github.com/gorilla/mux"
)

// Load in the templates we need
var templates = template.Must(template.ParseFiles("templates/home.html"))

func main() {
	r := mux.NewRouter()

	// This handles static media
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static")),
	))

	r.HandleFunc("/", home)
	r.HandleFunc("/upload", upload)

	http.ListenAndServe(":3000", r)
}

/******************************************************************************
 * Loads the home page where the user is greeted and can upload a file
 *****************************************************************************/
func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(5242880)

		// Load the uploaded file
    file, handler, err := r.FormFile("uploadfile")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

		// Create the temp file
    f, err := os.OpenFile("./temp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()

		// Copy the uploaded file to the temp file
    io.Copy(f, file)

	} else { // Regular get method
		// Get the user
		u, err := user.Current()
		if err != nil {
			fmt.Println(err)
		}

		// Load up the home template
		err = templates.ExecuteTemplate(w, "home.html", u.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

/******************************************************************************
 * Loads the home page where the user is greeted and can upload a file
 *****************************************************************************/
func upload(w http.ResponseWriter, r *http.Request) {

}
