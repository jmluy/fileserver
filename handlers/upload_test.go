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
    "io/ioutil"
    "log"
)

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func TestFileHandlerUpload(t *testing.T) {
    basePath := "/tmp/upload_test"
    subPath := "/myfiles"
    filename := "test-image.png"

    // create test directory in /tmp
    err := os.MkdirAll(filepath.Dir(basePath + subPath), os.ModePerm)
    checkErr(err)

    // open test file for upload
    file, err := os.Open(filename)
    checkErr(err)
    defer file.Close()

    bodyBuf := bytes.Buffer{}
    bodyWriter := multipart.NewWriter(&bodyBuf)
    fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
    checkErr(err)
    _, err = io.Copy(fileWriter, file)
    checkErr(err)
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

    fullPath := path.Clean(fmt.Sprintf("%v/%v/%v", basePath, subPath, filename))
    _, err = os.Open(fullPath)
    if err != nil {
        t.Errorf("File not uploaded", err)
    }

    // cleanup
    err = os.Remove(fullPath)
    checkErr(err)
}

func TestFileHandlerGet(t *testing.T) {
    basePath := "/tmp/upload_test"
    subPath := "/myfiles"
    filename := "test-image.png"

    data, err := ioutil.ReadFile(filename)
    checkErr(err)
    out := path.Clean(fmt.Sprintf("%v/%v/%v", basePath, subPath, filename))
    err = ioutil.WriteFile(out, data, 0644)
    checkErr(err)

    url := fmt.Sprintf("http://localhost:9898/file?path=%v&filename=%v", subPath, filename)
    handler := FileHandler(basePath)
    req, _ := http.NewRequest("GET", url, nil)

    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("Didn't receive a successful status", http.StatusOK)
    }
}
