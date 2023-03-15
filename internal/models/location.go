package models

import (
	"errors"
)

type Location struct {
	Id        int     `db:"id" json:"id"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
}

func (model *Location) GetLocationService() error {
	err := db.Get(model, "SELECT longitude, latitude FROM locations WHERE id=$1", model.Id)
	if err != nil {
		return errors.New("location not found")
	}
	return nil
}

func (model *Location) CreateLocationService() error {
	err := db.QueryRow("INSERT INTO locations (latitude, longitude) VALUES ($1, $2) RETURNING id", model.Latitude, model.Longitude).Scan(&model.Id)

	if err != nil {
		return errors.New("location already exist")
	}

	return nil
}

func (model *Location) UpdateLocationService() error {
	err := db.QueryRow("UPDATE locations SET latitude = $1, longitude = $2 WHERE id = $3 RETURNING id", model.Latitude, model.Longitude, model.Id).Scan(&model.Id)

	if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"uq_locations\"" {
		return errors.New("location already exist")
	} else if err != nil {
		return errors.New("location not found")
	}

	return nil
}

func (model *Location) DeleteLocationService() error {
	err := db.QueryRow("DELETE FROM locations WHERE id = $1 RETURNING id", model.Id).Scan(&model.Id)

	if err != nil {
		return errors.New("location not found")
	}

	return nil
}
