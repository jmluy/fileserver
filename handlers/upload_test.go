package handlers

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "bytes"
    "mime/multipart"
    "os"
    "io"
    "path"
    "fmt"
    "path/filepath"
)

func TestFileHandler(t *testing.T) {
    basePath := "/tmp/upload_test"
    subPath := "/myfiles"
    filename := "test-image.png"

    // create test directory in /tmp
    err := os.MkdirAll(filepath.Dir(basePath + subPath), os.ModePerm)
    if err != nil {
        panic(err)
    }

    // open test file for upload
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    bodyBuf := bytes.Buffer{}
    bodyWriter := multipart.NewWriter(&bodyBuf)
    fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
    if err != nil {
        panic(err)
    }
    _, err = io.Copy(fileWriter, file)
    if err != nil {
        panic(err)
    }
    bodyWriter.WriteField("path", subPath)

    bodyWriter.Close()

    handler := FileHandler(basePath)
    req, _ := http.NewRequest("POST", "http://localhost:9898/file", &bodyBuf)
    req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("Didn't receive a successful status", http.StatusOK)
    }

    fullPath := fmt.Sprintf("%v/%v/%v", basePath, subPath, filename)
    _, err = os.Open(path.Clean(fullPath))
    if err != nil {
        t.Errorf("File not uploaded", err)
    }
}
