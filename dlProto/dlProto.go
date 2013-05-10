package main

import (
	"code.google.com/p/goconf/conf"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type configValues struct {
	host      string
	port      string
	uploadDir string
	templates string
}
var config configValues

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
	files, err := ioutil.ReadDir(config.uploadDir)
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
	fileName := config.uploadDir + "/" + req.URL.Path

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

func parseConfigFile(file string) (configValues, error) {

	c, err := conf.ReadConfigFile(file)
	if err != nil {
		return config, err
	}

	config.host, err = c.GetString("server", "host")
	if err != nil {
		return config, err
	}

	config.port, err = c.GetString("server", "port")
	if err != nil {
		return config, err
	}

	config.uploadDir, err = c.GetString("directories", "uploadDir")
	if err != nil {
		return config, err
	}

	config.templates, err = c.GetString("directories", "templates")
	if err != nil {
		return config, err
	}

	return config, err
}

func parseTemplates() error {
	tmpl.index = template.Must(template.ParseFiles(
		config.templates + "/_base.html",
		config.templates + "/index.html",
	))

	tmpl.files = template.Must(template.ParseFiles(
		config.templates + "/_base.html",
		config.templates + "/files.html",
	))

	tmpl.upload = template.Must(template.ParseFiles(
		config.templates + "/_base.html",
		config.templates + "/upload.html",
	))

	tmpl.thanks = template.Must(template.ParseFiles(
		config.templates + "/_base.html",
		config.templates + "/thanks.html",
	))
	return nil
}

func main() {

	config, err := parseConfigFile("../webserver.conf")
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

	fmt.Printf("Webserver started %s:%s\n", config.host, config.port)
	if err := http.ListenAndServe(config.host+":"+config.port, nil); err != nil {
		panic(err)
	}
}
