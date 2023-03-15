package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"goSimbir/internal/dto"
	"goSimbir/internal/models"
	"goSimbir/utils"
	"net/http"
	"regexp"
	"strconv"
)

func validateAnimal(animal models.Animal) error {
	if utils.FindDublicates(animal.AnimalTypes) == true {
		return errors.New("finded duplicates in AnimalTypes")
	}

	if len(animal.AnimalTypes) <= 0 {
		return errors.New("invalid value")
	}

	for _, animalType := range animal.AnimalTypes {
		if animalType <= 0 {
			return errors.New("invalid value")
		}
	}

	if utils.Contains(animal.Gender, []string{"MALE", "FEMALE", "OTHER"}) == false {
		return errors.New("invalid value")
	}

	if animal.Weight <= 0 || animal.Length <= 0 || animal.Height <= 0 || animal.ChipperId <= 0 || animal.ChippingLocationId <= 0 {
		return errors.New("invalid value")
	}

	return nil
}

func createAnimal(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	animal := models.Animal{}
	_ = json.NewDecoder(r.Body).Decode(&animal)
	err := validateAnimal(animal)
	if err != nil {
		switch err.Error() {
		case "invalid value":
			w.WriteHeader(http.StatusBadRequest)
			return
		case "finded duplicates in AnimalTypes":
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	err = animal.CreateAnimalService()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(animal)
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	animal := models.Animal{}
	animal.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
	if animal.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := animal.GetAnimalService()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(animal)
}

func searchAnimals(w http.ResponseWriter, r *http.Request) {
	animal := models.Animal{}

	filterFields := dto.AnimalFindFields{}
	filterFields.StartDateTime = r.URL.Query().Get("startDateTime")
	// check iso 8601 with nano seconds
	match, _ := regexp.Match(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{1,9}Z`, []byte(filterFields.StartDateTime))
	if match == false && r.URL.Query().Get("startDateTime") != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filterFields.EndDateTime = r.URL.Query().Get("endDateTime")
	match, _ = regexp.Match(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{1,9}Z`, []byte(filterFields.EndDateTime))
	if match == false && r.URL.Query().Get("endDateTime") != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filterFields.ChipperId, _ = strconv.Atoi(r.URL.Query().Get("chipperId"))
	if filterFields.ChipperId == 0 && r.URL.Query().Get("chipperId") != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filterFields.ChippingLocationId, _ = strconv.Atoi(r.URL.Query().Get("chippingLocationId"))
	if filterFields.ChippingLocationId == 0 && r.URL.Query().Get("chippingLocationId") != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filterFields.LifeStatus = r.URL.Query().Get("lifeStatus")
	filterFields.Gender = r.URL.Query().Get("gender")

	var err error
	filterFields.From, err = strconv.Atoi(r.URL.Query().Get("from"))
	if r.URL.Query().Get("from") == "" {
		filterFields.From = 0
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if filterFields.From < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filterFields.Size, err = strconv.Atoi(r.URL.Query().Get("size"))
	if (r.URL.Query().Get("size")) == "" {
		filterFields.Size = 10
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if filterFields.Size <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	animals, err := animal.FindAnimalService(filterFields)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(animals)
}

func updateAnimal(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	animal := models.Animal{}
	animal.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
	_ = json.NewDecoder(r.Body).Decode(&animal)
	if animal.Id <= 0 || animal.Weight <= 0 || animal.Length <= 0 || animal.Height <= 0 ||
		utils.Contains(animal.Gender, []string{"MALE", "FEMALE", "OTHER"}) == false ||
		utils.Contains(animal.LifeStatus, []string{"ALIVE", "DEAD"}) == false ||
		animal.ChipperId <= 0 ||
		animal.ChippingLocationId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := animal.UpdateAnimalService()
	if err != nil {
		switch err.Error() {
		case "invalid value":
			animal.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
			w.WriteHeader(http.StatusBadRequest)
			return
		case "entity not found":
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(animal)
}

func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	animal := models.Animal{}
	animal.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])

	if animal.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := animal.DeleteAnimalService()
	if err != nil {
		switch err.Error() {
		case "invalid value":
			w.WriteHeader(http.StatusBadRequest)
			return
		case "animal not found":
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println(animal)
	_ = json.NewEncoder(w).Encode(animal)
}

func addAnimalTypeToAnimal(w http.ResponseWriter, r *http.Request) {
	animal := models.Animal{}
	animal.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
	animalTypeId, _ := strconv.Atoi(mux.Vars(r)["typeId"])
	if animal.Id <= 0 || animalTypeId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	err := animal.AddAnimalTypeToAnimalService(animalTypeId)
	if err != nil {
		switch err.Error() {
		case "typeId already exist for this animalId":
			w.WriteHeader(http.StatusConflict)
			return
		case "entity not found":
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(animal)
}

func updateAnimalTypeToAnimal(w http.ResponseWriter, r *http.Request) {
	animal := models.Animal{}
	animalEdit := dto.AnimalEdit{}
	animal.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
	_ = json.NewDecoder(r.Body).Decode(&animalEdit)
	if animal.Id <= 0 || animalEdit.OldTypeId <= 0 || animalEdit.NewTypeId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := animal.UpdateAnimalTypeToAnimalService(animalEdit)
	if err != nil {
		switch err.Error() {
		case "entity not found":
			w.WriteHeader(http.StatusNotFound)
			return
		case "already exist":
			w.WriteHeader(http.StatusConflict)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(animal)
}

func deleteAnimalTypeToAnimal(w http.ResponseWriter, r *http.Request) {
	animal := models.Animal{}
	animal.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
	var typeId int
	typeId, _ = strconv.Atoi(mux.Vars(r)["typeId"])
	if animal.Id <= 0 || typeId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	err := animal.DeleteAnimalTypeToAnimal(typeId)
	if err != nil {
		switch err.Error() {
		case "invalid value":
			w.WriteHeader(http.StatusBadRequest)
			return
		case "entity not found":
			w.WriteHeader(http.StatusNotFound)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(animal)
}
