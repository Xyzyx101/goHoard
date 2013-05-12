package main

import (
	"fmt"
	"github.com/Xyzyx101/goHoard/config"
	"io/ioutil"
	"log"
	"net/http"
)

var conf config.Config
var uri string

func dbExists() {
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err.Error(), "\nCouchDB not found.  Check couch is running and server.conf is correct")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		out, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(out), "\nDB already exists!")
	}
}

func createDB() {
	req, err := http.NewRequest("PUT", uri, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.SetBasicAuth(conf["couch"]["user"], conf["couch"]["pw"])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(resp.Status)
	out, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(out))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Database %s created successfully at %s:%s",
			conf["couch"]["db"],
			conf["couch"]["host"],
			conf["couch"]["port"])
	} else {
		log.Fatal("Database not created!")
	}
}

func main() {
	var err error
	conf, err = config.ParseFile("../server.conf")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	uri = "http://" + conf["couch"]["host"] + ":" + conf["couch"]["port"] + "/" + conf["couch"]["db"]

	dbExists()

	createDB()

	//	uploadDesignDoc()
}
