package server

import (
	"net/http"
	"seetime/server/account"
)

func openAPi() {
	http.HandleFunc("/api/login", account.HandleLogin)

	http.HandleFunc("/api/me", account.HandleMe)
	http.HandleFunc("/api/me/change", account.HandleMeChange)

	http.HandleFunc("/api/user", account.HandleUser)
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		// ---日志
		return
	}
}

func Loading() {
	defer openAPi()
}

func init() {
	SendInfo()
}
