package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/glog"
)

/*
PostResponder responds to a POST request with basic auth
*/
func PostResponder(w http.ResponseWriter, r *http.Request) {

	// Decode POST content
	decoder := json.NewDecoder(r.Body)

	var content request
	if err := decoder.Decode(&content); err != nil {
		panic(err)
	}

	// Basic auth
	usr, psw, ok := r.BasicAuth()
	glog.Infof("username: %s, password: %s, ok %t", usr, psw, ok)

	var (
		firstName string
		lastName  string
	)

	row := db.QueryRow(
		fmt.Sprintf(
			`SELECT SPRIDEN_LAST_NAME, SPRIDEN_FIRST_NAME
				FROM SATURN.SPRIDEN
				WHERE SPRIDEN_ID = '%s'`,
			strconv.Itoa(content.ID)))

	if err := row.Scan(&firstName, &lastName); err != nil {
		glog.Error(err)
	}

	response, _ := json.Marshal(
		&response{
			ID:        content.ID,
			FirstName: firstName,
			LastName:  lastName,
		})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}
