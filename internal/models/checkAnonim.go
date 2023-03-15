package models

import (
	b64 "encoding/base64"
	"net/http"
	"strings"
)

func CheckAnonim(r *http.Request) bool {

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		auth := strings.Split(authHeader, " ")[1]
		authDec, _ := b64.StdEncoding.DecodeString(auth)
		authDecArr := strings.Split(string(authDec), ":")
		login := authDecArr[0]
		password := authDecArr[1]
		// check login and password in db
		var dbPassword string
		err := db.QueryRow("SELECT password FROM accounts WHERE email = $1", login).Scan(&dbPassword)
		if err != nil || dbPassword != password {
			return true
		}
		return false
	}
	return true

}
