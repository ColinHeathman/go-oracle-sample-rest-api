package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Request struct
type request struct {
	Passthrough string `json:"passthrough"`
}

// Response struct
type response struct {
	Passthrough string `json:"passthrough"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

/*
PostResponder responds to a POST request with basic auth
*/
func PostResponder(w http.ResponseWriter, r *http.Request) {

	// Decode POST content
	var content request

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&content); err != nil {
		errorResponder(w, err)
		return
	}

	// Basic auth
	username, password, ok := r.BasicAuth()
	if !ok {
		// Unauthorized
		unauthorizedResponder(w)
		return
	}

	// Validate user with LDAP
	validated, err := Validate(username, password)
	if err != nil {
		errorResponder(w, err)
		return
	}
	if !validated {
		// Unauthorized
		unauthorizedResponder(w)
		return
	}

	// Get UDCID
	udcid, err := GetUDCID(username)
	if err != nil {
		errorResponder(w, err)
		return
	}

	// Database query
	var (
		firstName string
		lastName  string
	)

	// Query database
	row := db.QueryRow(
		fmt.Sprintf(
			`SELECT SPRIDEN_LAST_NAME, SPRIDEN_FIRST_NAME
				FROM SATURN.SPRIDEN
				WHERE SPRIDEN_ID = '%s'`,
			udcid))

	// Database row
	if err := row.Scan(&firstName, &lastName); err != nil {
		errorResponder(w, err)
		return
	}

	// Encode response
	response, err := json.Marshal(
		&response{
			Passthrough: content.Passthrough,
			FirstName:   firstName,
			LastName:    lastName,
		})
	if err != nil {
		errorResponder(w, err)
		return
	}

	// Write response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func unauthorizedResponder(w http.ResponseWriter) {
	// Write response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write([]byte(
		`{"message":"unauthorized"}`))
}

func errorResponder(w http.ResponseWriter, err error) {
	// Write response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(500)
	w.Write([]byte(
		fmt.Sprintf(`{"message":"%v"}`, err)))
}
