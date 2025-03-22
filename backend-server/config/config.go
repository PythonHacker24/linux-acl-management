package config 

import (
    "time"
    "sync"

    "backend-server/models"
)

const(
    Host = "localhost"
    Port = "9999"

    LdapServer = "ldap://localhost:389"
    LdapBindDN = "cn=admin,dc=example,dc=org"
    LdapPassword = "admin"
    LdapSearchBase = "dc=example,dc=org"
)

var Sessions = make(map[string]*models.Session)
var SessionMutex = sync.Mutex{}

const SessionTimeout = 5 * time.Minute
