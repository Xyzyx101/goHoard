package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var fileLocation = "/home/andrew/code/tuig/files/"

var indexPage = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/index.html",
))

var filesPage = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/files.html",
))

var uploadPage = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/upload.html",
))

func index(w http.ResponseWriter, req *http.Request) {
	if err := indexPage.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func filesHandler(w http.ResponseWriter, req *http.Request) {
	files := make([]os.FileInfo, 100)
	files, err := ioutil.ReadDir(fileLocation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	type FileList struct {
		Files []os.FileInfo
	}
	fileList := FileList{Files:files}
	if err := filesPage.Execute(w, fileList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ServeFileHandler struct {}
func (fh ServeFileHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fileName := fileLocation + req.URL.Path
	
	fileBuffer, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := w.Write(fileBuffer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func upload(w http.ResponseWriter, req *http.Request) {
	if err :=uploadPage.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func submitUpload(w http.ResponseWriter, req *http.Request) {
	formFile, fileHeader, err := req.FormFile("fileName")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fileName := fileLocation + fileHeader.Filename
	
	var filePerm os.FileMode = 0664
	
	fileBuffer, err :=ioutil.ReadAll(formFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	if err := ioutil.WriteFile(fileName, fileBuffer, filePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func main() {
	
	http.HandleFunc("/", index)
	
	http.HandleFunc("/files/", filesHandler)

	var getFile ServeFileHandler
	http.Handle("/file/", http.StripPrefix("/file/", getFile))

	http.HandleFunc("/upload/", upload)

	http.HandleFunc("/uploadFile", submitUpload)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
