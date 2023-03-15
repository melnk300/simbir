package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"goSimbir/internal/dto"
	"goSimbir/internal/models"
	"net/http"
	"regexp"
	"strconv"
)

func addVisitedLocation(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	visitedLocation := models.VisitedLocation{}
	visitedLocation.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
	visitedLocation.LocationPointId, _ = strconv.Atoi(mux.Vars(r)["pointId"])
	if visitedLocation.Id <= 0 || visitedLocation.LocationPointId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := visitedLocation.AddVisitedLocationService()
	if err != nil {
		switch err.Error() {
		case "not found":
			w.WriteHeader(http.StatusNotFound)
			return
		case "invalid value":
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(visitedLocation)
}

func updateVisitedLocation(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	visitedLocation := models.VisitedLocation{}
	visitedLocation.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
	_ = json.NewDecoder(r.Body).Decode(&visitedLocation)
	if visitedLocation.Id <= 0 || visitedLocation.LocationPointId <= 0 || visitedLocation.VisitedLocationPointId <= 0 || visitedLocation.LocationPointId == visitedLocation.VisitedLocationPointId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := visitedLocation.UpdateVisitedLocation()
	if err != nil {
		switch err.Error() {
		case "entity not found":
			w.WriteHeader(http.StatusNotFound)
			return
		case "invalid value":
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	visitedLocation.VisitedLocationPointId = 0
	_ = json.NewEncoder(w).Encode(visitedLocation)
}

func getVisitedLocation(w http.ResponseWriter, r *http.Request) {
	var err error
	visitedLocation := models.VisitedLocation{}
	filterFields := dto.VisitedLocationsFindFields{}
	visitedLocation.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])

	filterFields.StartDateTime = r.URL.Query().Get("startDateTime")
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

	filterFields.From, err = strconv.Atoi(r.URL.Query().Get("from"))
	if r.URL.Query().Get("from") == "" {
		filterFields.From = 0
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filterFields.Size, err = strconv.Atoi(r.URL.Query().Get("size"))
	if (r.URL.Query().Get("size")) == "" {
		filterFields.Size = 10
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if visitedLocation.Id <= 0 || filterFields.From < 0 || filterFields.Size <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	visitedLocations, err := visitedLocation.GetVisitedLocationService(filterFields)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(visitedLocations)
}

func deleteVisitedLocation(w http.ResponseWriter, r *http.Request) {
	if models.CheckAnonim(r) == true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	visitedLocation := models.VisitedLocation{}
	visitedLocation.Id, _ = strconv.Atoi(mux.Vars(r)["animalId"])
	visitedLocation.LocationPointId, _ = strconv.Atoi(mux.Vars(r)["visitedPointId"])
	if visitedLocation.Id <= 0 || visitedLocation.LocationPointId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := visitedLocation.DeleteVisitedLocationService()
	if err != nil {
		switch err.Error() {
		case "entity not found":
			w.WriteHeader(http.StatusNotFound)
			return
		case "invalid value":
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(visitedLocation)
}
