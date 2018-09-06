package main

import (
	"errors"
	"fmt"
	"log"

	ldap "gopkg.in/ldap.v2"
)

var (
	bindUsername = "cn=read-only-admin,dc=example,dc=com"
	bindPassword = "password"
	ldapHost     = "ldap.forumsys.com"
	baseDN       = "dc=example,dc=com"
	uidVariable  = "uid"
	// ldapHost     = "ldaplb.uni.adr.unbc.ca"
	// baseDN       = "OU=ITS,OU=Employees,OU=Users,OU=Core,DC=uni,DC=adr,DC=unbc,DC=ca"
	// uidVariable = "sAMAccountName"
)

/*
ConnectLDAP connects to the LDAP server
*/
func ConnectLDAP() *ldap.Conn {

	// Dial ldaps
	// l, err := ldap.DialTLS(
	// 	"tcp",
	// 	fmt.Sprintf("%s:%d", ldapHost, 636),
	// 	&tls.Config{
	// 		InsecureSkipVerify: true})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Dial ldap
	l, err := ldap.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", ldapHost, 389))
	if err != nil {
		log.Fatal(err)
	}

	// bind with a read only user
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		log.Fatal(err)
	}

	return l
}

/*
Validate validates a username and password
*/
func Validate(username string, password string) (validated bool, err error) {

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(%s=%s)", uidVariable, username),
		[]string{"dn"},
		nil,
	)

	sr, err := ad.Search(searchRequest)
	if err != nil {
		return false, err
	}

	if len(sr.Entries) != 1 {
		return false, nil
	}

	userdn := sr.Entries[0].DN

	// Dial ldaps
	// l, err := ldap.DialTLS(
	// 	"tcp",
	// 	fmt.Sprintf("%s:%d", ldapHost, 636),
	// 	&tls.Config{
	// 		InsecureSkipVerify: true})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Dial ldap
	l, err := ldap.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", ldapHost, 389))
	if err != nil {
		log.Fatal(err)
	}

	bindRequest := ldap.NewSimpleBindRequest(
		userdn,
		password,
		nil)

	result, err := l.SimpleBind(bindRequest)
	if err != nil {
		return false, nil
	}

	defer l.Close()

	return result != nil, nil
}

/*
GetUDCID gets udcid of a username
*/
func GetUDCID(username string) (udcid string, err error) {

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(%s=%s)", uidVariable, username),
		[]string{"dn", "cn", "mail"},
		nil,
	)

	sr, err := ad.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) != 1 {
		return "", errors.New("User not found")
	}

	return sr.Entries[0].GetAttributeValue("mail"), nil
}
