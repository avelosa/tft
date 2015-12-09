package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"io"
	"math/rand"
	"time"

	"github.com/gorilla/mux"
)

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
	// Load the templates
	t,_ := template.ParseFiles("templates/header.html", "templates/home.html",
		 "templates/footer.html")

	if r.Method == "GET" {
		t.ExecuteTemplate(w, "home.html", nil)
	} else if r.Method == "POST" { // Pressed upload
		r.ParseMultipartForm(5000000)

		// Load the uploaded file
		upload,_,err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			t.ExecuteTemplate(w, "home.html", "Error uploading the file")
			return
		}
		defer upload.Close()

		// Generate a random filename/url
		randString := RandomString(100)
		fmt.Println("random string =", randString)

		// Create the temp file
		temp, err := os.OpenFile("./temp/" +randString, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0666)
		if err != nil {
			fmt.Println(err)
			t.ExecuteTemplate(w, "home.html", "Error uploading the file")
			return
		}
		defer temp.Close()

		// Copy the uploaded file to the temp file
		_,err = io.Copy(temp, upload)
		if err != nil {
			fmt.Println(err)
			t.ExecuteTemplate(w, "home.html", "Error uploading the file")
			return
		}
		upload.Close()
		temp.Close()


	}
}

/******************************************************************************
 * Loads the home page where the user is greeted and can upload a file
 *****************************************************************************/
func upload(w http.ResponseWriter, r *http.Request) {

}


const possibleChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
/******************************************************************************
 * Generates a random string of length n only consisting of the above chars
 *****************************************************************************/
func RandomString(n int) string {
		rand.Seed(time.Now().UnixNano())
    b := make([]byte, n)
		numPos := len(possibleChars)
    for i := range b {
        b[i] = possibleChars[rand.Intn(numPos)]
    }
    return string(b)
}
