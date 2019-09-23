package utils

import (
	"fmt"
	"log"

	"github.com/mmitton/ldap"
)

func LdapBind(addr,
	BaseDN,
	BindDN,
	BindPasswd,
	UserField,
	user,
	password string) (sucess bool, err error) {

	filter := "(" + UserField + "=" + user + ")"
	conn, err := ldap.Dial("tcp", addr)

	if err != nil {
		return false, fmt.Errorf("dial ldap fail: %s", err.Error())
	}

	defer conn.Close()
	if BindDN != "" {
		err = conn.Bind(BindDN, BindPasswd)
	}
	if err != nil {
		return false, fmt.Errorf("ldap Bind fail: %s", err.Error())
	}
	search := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		nil,
		nil)

	sr, err := conn.Search(search)

	if err != nil {
		return false, fmt.Errorf("ldap search fail: %s", err.Error())
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("ERROR:", err)
			sucess = false
		}
	}()
	err = conn.Bind(sr.Entries[0].DN, password)
	return err == nil, err
}

func Ldapsearch(addr,
	BaseDN,
	BindDN,
	BindPasswd,
	UserField,
	user string,
	Attributes []string) (map[string]string, error) {

	filter := "(" + UserField + "=" + user + ")"
	conn, err := ldap.Dial("tcp", addr)

	if err != nil {
		return nil, fmt.Errorf("dial ldap fail: %s", err.Error())
	}

	defer conn.Close()

	if BindDN != "" {
		err = conn.Bind(BindDN, BindPasswd)
	}
	if err != nil {
		return nil, fmt.Errorf("ldap Bind fail: %s", err.Error())
	}
	search := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		Attributes,
		nil)

	sr, err := conn.Search(search)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR:", err)
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("ldap search fail: %s", err.Error())
	}
	var UserAttributes map[string]string
	UserAttributes = make(map[string]string)

	userSn := sr.Entries[0].GetAttributeValue(Attributes[0])
	userMail := sr.Entries[0].GetAttributeValue(Attributes[1])
	userTel := sr.Entries[0].GetAttributeValue(Attributes[2])

	UserAttributes["sn"] = userSn
	UserAttributes["telephoneNumber"] = userTel
	UserAttributes["mail"] = userMail
	return UserAttributes, err
}
