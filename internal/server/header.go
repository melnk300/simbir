package server

import (
	"github.com/gorilla/mux"
)

func initRoute() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/registration", registerAccount).Methods("POST")
	r.HandleFunc("/accounts/{accountId:-?[0-9]+}", getAccountById).Methods("GET")
	r.HandleFunc("/accounts/search", searchAccounts).Methods("GET")
	return r
}
