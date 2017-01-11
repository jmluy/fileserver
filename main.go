package main

import (
    "net/http"
    "github.com/jmluy/fileserver/handlers"
)

func main() {
    http.HandleFunc("/file", handlers.FileHandler("./uploads"))
    http.ListenAndServe(":9099", nil)
}