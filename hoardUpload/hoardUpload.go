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
	upload *template.Template
}

var tmpl templates

func postFile(w http.ResponseWriter, req *http.Request) {
	formFile, fileHeader, err := req.FormFile("fileName")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fileName := conf["dir"]["upload"] + "/" + fileHeader.Filename

	var filePerm os.FileMode = 0664

	fileBuffer, err := ioutil.ReadAll(formFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := ioutil.WriteFile(fileName, fileBuffer, filePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	thanksPage := "http://" + conf["webserver"]["host"] + ":" + conf["webserver"]["port"] + "/thanks"
	http.Redirect(w, req, thanksPage, http.StatusSeeOther)
}

func parseTemplates() error {
	tmpl.upload = template.Must(template.ParseFiles(
		conf["dir"]["template"]+"/_base.html",
		conf["dir"]["template"]+"/upload.html",
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

	http.HandleFunc("/postFile", postFile)

	host := conf["uploadserver"]["host"]
	port := conf["uploadserver"]["port"]
	fmt.Printf("Upload server started %s:%s\n", host, port)
	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		panic(err)
	}

}
