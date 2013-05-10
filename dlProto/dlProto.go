package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"code.google.com/p/goconf/conf"
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

var thanksPage = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/thanks.html",
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

func thanks(w http.ResponseWriter, req *http.Request) {
	if err :=thanksPage.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type configValues struct {
	host string
	port string
}

func parseConfigFile(file string) (configValues, error) {
	var config configValues

	c, err := conf.ReadConfigFile(file); 
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
		
	return config, err
}

func main() {
	
	config, err := parseConfigFile("webserver.conf")
	if err != nil {
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
	if err := http.ListenAndServe(config.host + ":" + config.port, nil); err != nil {
		panic(err)
	}
}
