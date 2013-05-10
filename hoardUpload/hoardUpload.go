package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

var fileLocation = "/home/andrew/code/tuig/files/"

func postFile(w http.ResponseWriter, req *http.Request) {
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

	http.Redirect(w, req, "http://localhost:8080/thanks", http.StatusSeeOther)

}

func main() {

	http.HandleFunc("/postFile", postFile)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}

}
