package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"goSimbir/internal/models"
	"net/http"
	"strconv"
)

type LocationBody struct {
	Id        int      `db:"id" json:"id"`
	Latitude  *float64 `db:"latitude" json:"latitude"`
	Longitude *float64 `db:"longitude" json:"longitude"`
}

func getLocation(w http.ResponseWriter, r *http.Request) {
	location := models.Location{}
	location.Id, _ = strconv.Atoi(mux.Vars(r)["locationId"])
	if location.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := location.GetLocationService()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(location)
}

func createLocation(w http.ResponseWriter, r *http.Request) {
	rawLocation := LocationBody{}
	_ = json.NewDecoder(r.Body).Decode(&rawLocation)

	if rawLocation.Latitude == nil || rawLocation.Longitude == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	location := models.Location{Longitude: *rawLocation.Longitude, Latitude: *rawLocation.Latitude}

	if location.Latitude > 90 || location.Latitude < -90 || location.Longitude > 180 || location.Longitude < -180 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := location.CreateLocationService()
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(location)
}

func updateLocation(w http.ResponseWriter, r *http.Request) {
	rawLocation := LocationBody{}
	_ = json.NewDecoder(r.Body).Decode(&rawLocation)

	if rawLocation.Latitude == nil || rawLocation.Longitude == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	location := models.Location{Longitude: *rawLocation.Longitude, Latitude: *rawLocation.Latitude}
	location.Id, _ = strconv.Atoi(mux.Vars(r)["locationId"])

	if location.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if location.Latitude > 90 || location.Latitude < -90 || location.Longitude > 180 || location.Longitude < -180 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := location.UpdateLocationService()
	if err != nil {
		switch err.Error() {
		case "location not found":
			w.WriteHeader(http.StatusNotFound)
			return
		case "location already exist":
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(location)
}
