package main

import (
	"fmt"
	"github.com/Xyzyx101/goHoard/config"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var conf config.Config

type templates struct {
	index  *template.Template
	upload *template.Template
	files  *template.Template
	thanks *template.Template
}

var tmpl templates

func index(w http.ResponseWriter, req *http.Request) {
	if err := tmpl.index.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func filesHandler(w http.ResponseWriter, req *http.Request) {
	files := make([]os.FileInfo, 100)
	files, err := ioutil.ReadDir(conf["dir"]["upload"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	type FileList struct {
		Files []os.FileInfo
	}
	fileList := FileList{Files: files}
	if err := tmpl.files.Execute(w, fileList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ServeFileHandler struct{}

func (fh ServeFileHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fileName := conf["dir"]["upload"] + "/" + req.URL.Path

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
	if err := tmpl.upload.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func thanks(w http.ResponseWriter, req *http.Request) {
	if err := tmpl.thanks.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseTemplates() error {
	tmpl.index = template.Must(template.ParseFiles(
		conf["dir"]["template"]+"/_base.html",
		conf["dir"]["template"]+"/index.html",
	))

	tmpl.files = template.Must(template.ParseFiles(
		conf["dir"]["template"]+"/_base.html",
		conf["dir"]["template"]+"/files.html",
	))

	tmpl.upload = template.Must(template.ParseFiles(
		conf["dir"]["template"]+"/_base.html",
		conf["dir"]["template"]+"/upload.html",
	))

	tmpl.thanks = template.Must(template.ParseFiles(
		conf["dir"]["template"]+"/_base.html",
		conf["dir"]["template"]+"/thanks.html",
	))
	return nil
}

func main() {
	var err error
	conf, err = config.ParseFile("../server.conf")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := parseTemplates(); err != nil {
		fmt.Println(err.Error())
		return
	}

	http.HandleFunc("/", index)

	http.HandleFunc("/files/", filesHandler)

	var getFile ServeFileHandler
	http.Handle("/file/", http.StripPrefix("/file/", getFile))

	http.HandleFunc("/upload/", upload)

	http.HandleFunc("/thanks/", thanks)

	fmt.Printf("Webserver started %s:%s\n", conf["webserver"]["host"], conf["webserver"]["port"])
	if err := http.ListenAndServe(conf["webserver"]["host"]+":"+conf["webserver"]["port"], nil); err != nil {
		panic(err)
	}
}
