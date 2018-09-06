package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/ldap.v2"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"gopkg.in/cas.v2"
)

type request struct {
	ID int `json:"id"`
}

type response struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var db *sql.DB
var ad *ldap.Conn

func main() {
	flag.Parse()

	glog.Info("Starting up")

	glog.Info("Connecting to database")

	db = ConnectDB()
	ad = ConnectLDAP()

	casURL, _ := url.Parse("https://cas.unbc.ca/cas")
	client := cas.NewClient(&cas.Options{
		URL: casURL,
	})

	r := mux.NewRouter()

	r.HandleFunc("/health", home).Methods("GET")
	r.HandleFunc("/get", GetResponder).Methods("GET")
	r.HandleFunc("/post", PostResponder).Methods("POST")

	if err := http.ListenAndServe(":8080", client.Handle(r)); err != nil {
		glog.Errorf("Error from HTTP Server: %v", err)
	}

	glog.Info("Shutting down")

	defer ad.Close()
	defer db.Close()

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok") 
}
