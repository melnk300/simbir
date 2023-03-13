package server

import (
	"github.com/gorilla/mux"
)

func initRoute() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/registration", registerAccount).Methods("POST")
	r.HandleFunc("/accounts/{accountId:-?[0-9]+}", getAccount).Methods("GET")
	r.HandleFunc("/accounts/search", searchAccounts).Methods("GET")
	r.HandleFunc("/accounts/{accountId:-?[0-9]+}", updateAccount).Methods("PUT")
	r.HandleFunc("/accounts/{accountId:-?[0-9]+}", deleteAccount).Methods("DELETE")

	r.HandleFunc("/locations/{locationId:-?[0-9]+}", getLocation).Methods("GET")
	r.HandleFunc("/locations", createLocation).Methods("POST")
	r.HandleFunc("/locations/{locationId:-?[0-9]+}", updateLocation).Methods("PUT")
	return r
}
