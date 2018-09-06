package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	cas "gopkg.in/cas.v2"
	ldap "gopkg.in/ldap.v2"
)

var db *sql.DB
var ad *ldap.Conn

func main() {
	// Parse CLI flags
	flag.Parse()

	glog.Info("Starting up")

	glog.Info("Connecting to database")
	db = ConnectDB()

	// Query database
	row := db.QueryRow(`SELECT 1 FROM DUAL;`)
	var dummy string
	if err := row.Scan(&dummy); err != nil {
		log.Fatal(err)
	}
	if dummy != "1" {
		log.Fatal("Failed to connect to Database")
	}

	glog.Info("Connecting to LDAP")
	ad = ConnectLDAP()

	// CAS configuration
	casURL, _ := url.Parse("https://cas.unbc.ca/cas")
	client := cas.NewClient(&cas.Options{
		URL: casURL,
	})

	// Request path router
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

	// Query database
	row := db.QueryRow(`SELECT 1 FROM DUAL;`)
	var dummy string
	if err := row.Scan(&dummy); err != nil {
		log.Fatal(err)
	}
	if dummy != "1" {
		fmt.Fprintf(w, "Failed to connect to Database")
	} else {
		fmt.Fprintf(w, "ok")
	}

}
