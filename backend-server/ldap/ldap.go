package ldap

import (
	"backend-server/config"
	"fmt"
	"log/slog"

	"github.com/go-ldap/ldap/v3"
)

// For failed authentications, we will log the requests and not for authenticated ones.
// Those who are properly authenticated don't need to be in the logs, they will just be printed on the console. 
// Failed authentication can potentially mean attempt to unauthorized access, so it would be logged.

func AuthenticateUser(username, userPassword, searchBase string) bool {
    l, err := ldap.DialURL(config.LdapServer)
    if err != nil {
        slog.Error("Failed to connect to LDAP Server", "Error", err.Error())
    }
    defer l.Close()

    err = l.Bind(config.LdapBindDN, config.LdapPassword)
    if err != nil {
        slog.Error("Admin authentication failed", "Error", err.Error())
    }

    searchRequest := ldap.NewSearchRequest(
		searchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(uid=%s)", username), // Searching by username
		[]string{"dn"}, // We only need the DN
		nil,
	)

    sr, err := l.Search(searchRequest)
	if err != nil || len(sr.Entries) == 0 {
		slog.Error("User not found in LDAP", "Error", err.Error())
		return false
	}

    userDN := sr.Entries[0].DN 

    err = l.Bind(userDN, userPassword)
	if err != nil {
		slog.Error("User authentication failed", "Error", err.Error())
		return false
	}

    fmt.Println("User authentication successful")
	return true
}
