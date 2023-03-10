package server

import (
	"github.com/gorilla/mux"
)

func initRoute() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/registration", registerAccount).Methods("POST")

	return r
}
