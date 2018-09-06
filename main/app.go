package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	cas "gopkg.in/cas.v2"
	ldap "gopkg.in/ldap.v2"
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

	router := mux.NewRouter()

	router.HandleFunc("/health", home).Methods("GET")
	router.HandleFunc("/get", GetResponder).Methods("GET")
	router.HandleFunc("/post", PostResponder).Methods("POST")

	if err := http.ListenAndServe(":8080", client.Handle(router)); err != nil {
		glog.Errorf("Error from HTTP Server: %v", err)
	}

	glog.Info("Shutting down")

	defer ad.Close()
	defer db.Close()

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
