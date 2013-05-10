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
	host          string
	port          string
	uploadDir     string
	templates     string
	webserverHost string
	webserverPort string
}

var config configValues

type templates struct {
	upload *template.Template
}

var tmpl templates

func postFile(w http.ResponseWriter, req *http.Request) {
	formFile, fileHeader, err := req.FormFile("fileName")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fileName := config.uploadDir + "/" + fileHeader.Filename

	var filePerm os.FileMode = 0664

	fileBuffer, err := ioutil.ReadAll(formFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := ioutil.WriteFile(fileName, fileBuffer, filePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	thanksPage := "http://" + config.webserverHost + ":" + config.webserverPort + "/thanks"
	http.Redirect(w, req, thanksPage, http.StatusSeeOther)
}

func parseConfigFile(file string) (configValues, error) {

	c, err := conf.ReadConfigFile(file)
	if err != nil {
		return config, err
	}

	config.host, err = c.GetString("uploadserver", "host")
	if err != nil {
		return config, err
	}

	config.port, err = c.GetString("uploadserver", "port")
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
	config.webserverHost, err = c.GetString("webserver", "host")
	if err != nil {
		return config, err
	}

	config.webserverPort, err = c.GetString("webserver", "port")
	if err != nil {
		return config, err
	}

	return config, err
}

func parseTemplates() error {
	tmpl.upload = template.Must(template.ParseFiles(
		config.templates+"/_base.html",
		config.templates+"/upload.html",
	))
	return nil
}

func main() {
	config, err := parseConfigFile("../server.conf")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := parseTemplates(); err != nil {
		fmt.Println(err.Error())
		return
	}

	http.HandleFunc("/postFile", postFile)

	fmt.Printf("Upload server started %s:%s\n", config.host, config.port)
	if err := http.ListenAndServe(config.host+":"+config.port, nil); err != nil {
		panic(err)
	}

}
