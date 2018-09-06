package main

import (
	"crypto/tls"
	"fmt"
	"log"

	ldap "gopkg.in/ldap.v2"
)

var (
	bindusername = "************"
	bindpassword = "************"
)

/*
ConnectLDAP connects to the LDAP server
*/
func ConnectLDAP() *ldap.Conn {

	// Dial ldap
	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", "ldaplb.uni.adr.unbc.ca", 636),
		&tls.Config{
			InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err)
	}

	// bind with a read only user
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}

	return l
}

/*
Validate validates a username and password
*/
func Validate(username string, password string) bool {

	// Dial ldap
	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", "ldaplb.uni.adr.unbc.ca", 636),
		&tls.Config{
			InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err)
	}

	bindRequest := ldap.NewSimpleBindRequest(
		fmt.Sprintf("CN=%s,OU=ITS,OU=Employees,OU=Users,OU=Core,DC=uni,DC=adr,DC=unbc,DC=ca", username),
		password,
		nil)

	r, err := l.SimpleBind(bindRequest)
	if err != nil {
		return false
	}

	defer l.Close()

	return r != nil
}

/*
GetUDCID gets udcid of a username
*/
func GetUDCID(username string) (udcid string, err error) {

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		"OU=ITS,OU=Employees,OU=Users,OU=Core,DC=uni,DC=adr,DC=unbc,DC=ca",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(sAMAccountName=%s)", username),
		[]string{"dn"},
		nil,
	)

	result, err := ad.Search(searchRequest)
	if err != nil {
		return "", err
	}
	attributes := result.Entries[0].Attributes[:]

	for i := 0; i < len(attributes); i++ {
		println(attributes[i])
	}

	return "", nil
}
