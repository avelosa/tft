package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os/user"

	"io"
	"time"
	"os"
	"strconv"
	"crypto/md5"

	"github.com/gorilla/mux"
)

func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method)
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        t, _ := template.ParseFiles("upload.gtpl")
        t.Execute(w, token)
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
        fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile("./uploaded/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
    }
}

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
	r.HandleFunc("/upload", upload)

	http.ListenAndServe(":3000", r)
}
