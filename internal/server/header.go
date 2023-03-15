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
	r.HandleFunc("/locations/{locationId:-?[0-9]+}", deleteLocation).Methods("DELETE")

	r.HandleFunc("/animals/types/{animalTypeId:-?[0-9]+}", getAnimalType).Methods("GET")
	r.HandleFunc("/animals/types/{animalTypeId:-?[0-9]+}", updateAnimalType).Methods("PUT")
	r.HandleFunc("/animals/types", createAnimalType).Methods("POST")
	r.HandleFunc("/animals/types/{animalTypeId:-?[0-9]+}", deleteAnimalType).Methods("DELETE")

	r.HandleFunc("/animals/{animalId:-?[0-9]+}", getAnimal).Methods("GET")
	r.HandleFunc("/animals/search", searchAnimals).Methods("GET")
	r.HandleFunc("/animals/{animalId:-?[0-9]+}", updateAnimal).Methods("PUT")
	r.HandleFunc("/animals", createAnimal).Methods("POST")
	r.HandleFunc("/animals/{animalId:-?[0-9]+}", deleteAnimal).Methods("DELETE")
	r.HandleFunc("/animals/{animalId:-?[0-9]+}/types/{typeId:-?[0-9]+}", addAnimalTypeToAnimal).Methods("POST")
	r.HandleFunc("/animals/{animalId:-?[0-9]+}/types", updateAnimalTypeToAnimal).Methods("PUT")
	r.HandleFunc("/animals/{animalId:-?[0-9]+}/types/{typeId:-?[0-9]+}", deleteAnimalTypeToAnimal).Methods("DELETE")

	r.HandleFunc("/animals/{animalId:-?[0-9]+}/locations/{pointId:-?[0-9]+}", addVisitedLocation).Methods("POST")
	r.HandleFunc("/animals/{animalId:-?[0-9]+}/locations", updateVisitedLocation).Methods("PUT")
	r.HandleFunc("/animals/{animalId:-?[0-9]+}/locations", getVisitedLocation).Methods("GET")
	r.HandleFunc("/animals/{animalId:-?[0-9]+}/locations/{visitedPointId:-?[0-9]+}", deleteVisitedLocation).Methods("DELETE")
	return r
}
