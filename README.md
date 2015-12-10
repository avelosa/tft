# tft
tft will be a simple temporary file transfer web app written in Go.

To build and run:
```
 go build app.go
 app
 ```
 
 To delete files, in current directory, older than 30 minutes. Should only be ran from in the /temp folder.
 ```
 find * -mtime +30m -type f -delete
 ```

## External Libraries
 Download and build:
 ```
 go get github.com/gorilla/mux
 ```
