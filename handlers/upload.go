package handlers

import (
    "net/http"
    "fmt"
    "os"
    "io"
    "path/filepath"
    "log"
)

const MB_SIZE = 1 << 20

var baseDirectory = "./uploads"

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(8 * MB_SIZE)
    //path := r.Form.Get("path")
    path := r.PostForm.Get("path")
    if path == "" {
        log.Println("No path specified")
        http.Error(w, "No path specified", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("uploadfile")
    if err != nil {
        http.Error(w, "Error uploading file", http.StatusBadRequest)
        return
    }

    if file == nil {
        http.Error(w, "No file specified", http.StatusBadRequest)
        return
    }

    defer file.Close()
    fmt.Println("Uploaded file info: ", handler.Header)
    fileWithPath := fmt.Sprintf("%v/%v", filepath.ToSlash(path), handler.Filename)
    localFilename := fmt.Sprintf(baseDirectory + "%v", filepath.Clean(fileWithPath))
    err = os.MkdirAll(filepath.Dir(localFilename), os.ModePerm)
    checkError(err)
    f, err := os.OpenFile(localFilename, os.O_WRONLY | os.O_CREATE, 0666)
    checkError(err)
    defer f.Close()
    _, err = io.Copy(f, file)
    checkError(err)

    w.Write([]byte("Successfully uploaded"))
}

func FileHandler(basePath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        baseDirectory = basePath

        if r.Method == "POST" || r.Method == "PUT" {
            uploadFile(w, r)
        } else {
            http.Error(w, "Not supported", http.StatusMethodNotAllowed)
        }
    }
}

