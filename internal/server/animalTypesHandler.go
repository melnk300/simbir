package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"goSimbir/internal/models"
	"net/http"
	"strconv"
)

func updateAnimalType(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	animalType := models.AnimalType{}

	animalType.Id, _ = strconv.Atoi(mux.Vars(r)["animalTypeId"])
	if animalType.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&animalType)
	if validateField(animalType.Type) == false {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := animalType.UpdateAnimalTypeService()
	if err != nil {
		switch err.Error() {
		case "type already created":
			w.WriteHeader(http.StatusConflict)
			return
		case "type not found":
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(animalType)
}

func createAnimalType(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	animalType := models.AnimalType{}
	_ = json.NewDecoder(r.Body).Decode(&animalType)
	if validateField(animalType.Type) {
		err := animalType.CreateAnimalTypeService()
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(animalType)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getAnimalType(w http.ResponseWriter, r *http.Request) {
	animalType := models.AnimalType{}
	animalType.Id, _ = strconv.Atoi(mux.Vars(r)["animalTypeId"])
	if animalType.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := animalType.GetAnimalTypeByIdService()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(animalType)
}

func deleteAnimalType(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	animalType := models.AnimalType{}
	animalType.Id, _ = strconv.Atoi(mux.Vars(r)["animalTypeId"])
	if animalType.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := animalType.DeleteAnimalTypeService()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
