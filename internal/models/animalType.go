package models

import (
	"errors"
)

type AnimalType struct {
	Id   int    `db:"id" json:"id"`
	Type string `db:"type" json:"type"`
}

func (model *AnimalType) UpdateAnimalTypeService() error {
	err := db.QueryRow("UPDATE animal_types SET type = $1 WHERE id = $2 RETURNING id", model.Type, model.Id).Scan(&model.Id)

	if err != nil {
		switch err.Error() {
		case "pq: duplicate key value violates unique constraint \"animal_types_type_key\"":
			return errors.New("type already created")
		case "sql: no rows in result set":
			return errors.New("type not found")
		}
	}

	return nil
}

func (model *AnimalType) CreateAnimalTypeService() error {
	err := db.QueryRow("INSERT INTO animal_types (type) VALUES ($1) RETURNING id", model.Type).Scan(&model.Id)
	if err != nil {
		return errors.New("type already created")
	}
	return nil
}

func (model *AnimalType) GetAnimalTypeByIdService() error {
	err := db.Get(model, "SELECT id, type FROM animal_types WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("account already created")
	}
	return nil
}

func (model *AnimalType) DeleteAnimalTypeService() error {
	err := db.QueryRow("DELETE FROM animal_types WHERE id = $1 RETURNING id", model.Id).Scan(&model.Id)
	if err != nil {
		return errors.New("type not found")
	}
	return nil
}
