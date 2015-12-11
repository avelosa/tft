package main

import (
	"archive/zip"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var hostname = "http://localhost:3000"

func init() {
	// Seed random
	rand.Seed(time.Now().UnixNano())
}

func main() {
	r := mux.NewRouter()

	// This handles static media
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static")),
	))

	r.HandleFunc("/", home)
	r.HandleFunc("/upload/{file}", upload)
	r.HandleFunc("/download/{file}", download)

	http.ListenAndServe(":3000", r)
}

type HomeTemplateData struct {
	Err     string
	Success string
}

type DownloadTemplateData struct {
	Filename string
}

/******************************************************************************
 * Loads the home page where the user is greeted and can upload a file
 *****************************************************************************/
func home(w http.ResponseWriter, r *http.Request) {
	// Load the templates
	t, _ := template.ParseFiles("templates/header.html", "templates/home.html",
		"templates/footer.html")

	if r.Method == "GET" {
		t.ExecuteTemplate(w, "home.html", HomeTemplateData{"", ""})
	} else if r.Method == "POST" { // Pressed upload
		r.ParseMultipartForm(5000000)

		// Generate a random archive file
		archive, err, name := RandomArchiveFile(100)
		if err != nil {
			fmt.Println(err)
			t.ExecuteTemplate(w, "home.html", HomeTemplateData{"Error uploading the file", ""})
			return
		}
		defer archive.Close()

		// TODO: support uploading of multiple files at once into the archive?
		// Load the uploaded file
		upload, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			t.ExecuteTemplate(w, "home.html", HomeTemplateData{"Error uploading the file", ""})
			return
		}
		defer upload.Close()

		// Make the file in the archive
		f, err := archive.Create(handler.Filename)
		if err != nil {
			fmt.Println(err)
			t.ExecuteTemplate(w, "home.html", HomeTemplateData{"Error uploading the file", ""})
			return
		}
		// Copy the data into the file
		_, err = io.Copy(f, upload)
		if err != nil {
			fmt.Println(err)
			t.ExecuteTemplate(w, "home.html", HomeTemplateData{"Error uploading the file", ""})
			return
		}
		upload.Close()

		// Close archive
		archive.Close()

		// Inform user of success and temp url
		t.ExecuteTemplate(w, "home.html", HomeTemplateData{"", hostname + "/download/" + name})
	}
}

/******************************************************************************
 * Used to access an uploaded file
 *****************************************************************************/
func upload(w http.ResponseWriter, r *http.Request) {
	// Get the variable filename
	filename := "./temp/" + mux.Vars(r)["file"] + ".zip"

	// See if the file exists
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		io.WriteString(w, "No files being hosted at this url")
		return
	}
	file.Close()

	// Download file
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=temp.zip")
	http.ServeFile(w, r, filename)
}

func download(w http.ResponseWriter, r *http.Request) {
	// Load the templates
	t, _ := template.ParseFiles("templates/header.html", "templates/download.html",
		"templates/footer.html")

	// Get the variable filename
	filename := "./temp/" + mux.Vars(r)["file"] + ".zip"

	// Inform user of success and temp url
	t.ExecuteTemplate(w, "download.html", DownloadTemplateData{hostname + "/download/" + filename})
}

const validFileChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

/******************************************************************************
 * Generates a random archive file in the temp/ folder
 *****************************************************************************/
func RandomArchiveFile(n int) (*zip.Writer, error, string) {
	b := make([]byte, n)
	numPos := len(validFileChars)

	// Generate random bytes
	for i := range b {
		b[i] = validFileChars[rand.Intn(numPos)]
	}
	// Convert bytes to string
	name := string(b)

	// Attempt to open the file
	file, err := os.OpenFile("./temp/"+name+".zip", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err, ""
	}

	// Return archived file
	return zip.NewWriter(file), nil, name
}
