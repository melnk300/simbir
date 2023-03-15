package models

import (
	"errors"
	"fmt"
	"goSimbir/internal/dto"
	"time"
)

type VisitedLocation struct {
	Id                           int       `db:"id" json:"id"`
	DateTimeOfVisitLocationPoint time.Time `db:"date_time_of_visited_location_point" json:"dateTimeOfVisitLocationPoint"`
	VisitedLocationPointId       int       `json:"visitedLocationPointId,omitempty" db:"location_point_id"`
	LocationPointId              int       `db:"location_point_id" json:"locationPointId"`
}

func (model *VisitedLocation) AddVisitedLocationService() error {
	animal := Animal{}
	err := db.Get(&animal, "SELECT id, chipping_location_id, life_status FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("not found")
	}

	if animal.LifeStatus == "DEAD" {
		return errors.New("invalid value")
	}

	var pointId int
	err = db.Get(&pointId, "SELECT id FROM locations WHERE id = $1", model.LocationPointId)
	if err != nil {
		return errors.New("not found")
	}

	var currentLocation int
	err = db.Get(&currentLocation, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1 ORDER BY id DESC LIMIT 1", model.Id)
	if err != nil {
		if animal.ChippingLocationId == model.LocationPointId {
			return errors.New("invalid value")
		}
	} else if currentLocation == model.LocationPointId {
		return errors.New("invalid value")
	}

	err = db.QueryRow("INSERT INTO \"animals-locations\"  (animal_id, location_id) VALUES ($1, $2) RETURNING visited_date_time",
		model.Id, model.LocationPointId).Scan(&model.DateTimeOfVisitLocationPoint)
	if err != nil {
		panic(err)
	}
	return nil
}

func (model *VisitedLocation) UpdateVisitedLocation() error {
	animal := Animal{}

	err := db.Get(&animal, "SELECT chipping_location_id FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	}
	var animalId int
	err = db.Get(&animalId, "SELECT id FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	}

	var checkOldLocationId int
	err = db.Get(&checkOldLocationId, "SELECT id FROM locations WHERE id = $1", model.VisitedLocationPointId)
	if err != nil {
		return errors.New("entity not found")
	}

	var checkNewLocationId int
	err = db.Get(&checkNewLocationId, "SELECT id FROM locations WHERE id = $1", model.LocationPointId)
	if err != nil {
		return errors.New("entity not found")
	}

	var firstVisitedLocation int
	err = db.Get(&firstVisitedLocation, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1 LIMIT 1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	} else if firstVisitedLocation == model.VisitedLocationPointId && model.VisitedLocationPointId == animal.ChippingLocationId {
		return errors.New("invalid value")
	}

	var existingVisitedLocationPointId int
	err = db.Get(&existingVisitedLocationPointId, "SELECT id FROM \"animals-locations\" WHERE location_id = $1 AND animal_id = $2 ORDER BY id DESC LIMIT 1",
		model.VisitedLocationPointId, model.Id)
	//if err != nil {
	//	return errors.New("not found")
	//}
	model.Id = existingVisitedLocationPointId

	var lowerLocationId int
	var lowerLocationRowId int
	err = db.QueryRow("SELECT location_id, id FROM \"animals-locations\" WHERE id < $1 AND animal_id=$2 ORDER BY id DESC LIMIT 1", existingVisitedLocationPointId, animalId).Scan(&lowerLocationId, &lowerLocationRowId)
	if err != nil {
		fmt.Println(err.Error())
	}

	var upperLocationId int
	var upperLocationRowId int
	err = db.QueryRow("SELECT location_id, id FROM \"animals-locations\" WHERE id > $1 AND animal_id=$2 ORDER BY id LIMIT 1", existingVisitedLocationPointId, animalId).Scan(&upperLocationId, &upperLocationRowId)
	if err != nil {
		fmt.Println(err.Error())
	}

	if lowerLocationId == model.LocationPointId || upperLocationId == model.LocationPointId {
		return errors.New("invalid value")
	}

	err = db.QueryRow("UPDATE \"animals-locations\" SET location_id = $1 WHERE id=$2 RETURNING id", model.LocationPointId, existingVisitedLocationPointId).Scan(&model.Id)
	if err != nil {
		return errors.New("entity not found")
	}

	err = db.QueryRow("SELECT visited_date_time FROM \"animals-locations\" WHERE id=$1", existingVisitedLocationPointId).Scan(&model.DateTimeOfVisitLocationPoint)
	if err != nil {
		panic(err)
	}
	return nil

	//animal := Animal{}
	//err := db.Get(&animal, "SELECT * FROM animals WHERE id = $1", model.Id)
	//if err != nil {
	//	// Животное с animalId не найдено
	//	return errors.New("entity not found")
	//}
	//
	//var firstVisitedLocationId int
	//err = db.Get(&firstVisitedLocationId, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1 LIMIT 1", model.Id)
	//if err == nil && firstVisitedLocationId == model.VisitedLocationPointId {
	//	// Обновление первой посещенной точки на точку чипирования
	//	return errors.New("invalid value")
	//}
	//
	//var lastOldLocationPointRowId int
	//err = db.Get(&lastOldLocationPointRowId, "SELECT id FROM \"animals-locations\" WHERE animal_id = $1 AND location_id = $2 ORDER BY id DESC LIMIT 1", model.Id, model.VisitedLocationPointId)
	//if err != nil {
	//	// У животного нет объекта с информацией о посещенной точке локации с visitedLocationPointId.
	//	return errors.New("entity not found")
	//} else {
	//	var lastOldLocationPointRowIdLower int
	//	err = db.Get(&lastOldLocationPointRowIdLower, "SELECT id FROM \"animals-locations\" WHERE animal_id = $1 AND id < $2 ORDER BY id DESC LIMIT 1", model.Id, lastOldLocationPointRowId)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//
	//	var lastOldLocationPointRowIdUpper int
	//	err = db.Get(&lastOldLocationPointRowIdUpper, "SELECT id FROM \"animals-locations\" WHERE animal_id = $1 AND id > $2 LIMIT 1", model.Id, lastOldLocationPointRowId)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//
	//	if model.LocationPointId == lastOldLocationPointRowIdLower || model.LocationPointId == lastOldLocationPointRowIdUpper {
	//		// Обновление точки локации на точку, совпадающую со следующей и/или с предыдущей точками
	//		return errors.New("invalid value")
	//	}
	//}
	//
	//var NewLocationPointRowId int
	//err = db.Get(&NewLocationPointRowId, "SELECT id FROM locations WHERE id = $1", model.LocationPointId)
	//if err != nil {
	//	// Точка локации с locationPointId не найдена
	//	return errors.New("entity not found")
	//}
	//
	//err = db.QueryRow("UPDATE \"animals-locations\" SET location_id = $1 WHERE id=$2 RETURNING id", model.LocationPointId, lastOldLocationPointRowId).Scan(&model.Id)
	//return nil
}

func (model *VisitedLocation) GetVisitedLocationService(fields dto.VisitedLocationsFindFields) ([]VisitedLocation, error) {
	var visitedLocations []VisitedLocation

	var animalIdCheck int //checking not found err for animal id
	err := db.Get(&animalIdCheck, "SELECT id FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return nil, errors.New("entity not found")
	}

	query := `SELECT animal_id, visited_date_time, location_id FROM "animals-locations"`
	if fields.StartDateTime != "" {
		query += fmt.Sprintf(" AND visited_date_time >= '%s'", fields.StartDateTime)
	}
	if fields.EndDateTime != "" {
		query += fmt.Sprintf(" AND visited_date_time <= '%s'", fields.EndDateTime)
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", fields.Size, fields.From)

	err = db.Select(&visitedLocations, query)
	if err != nil {
		return nil, errors.New("entity not found")
	}

	if len(visitedLocations) == 0 {
		return nil, errors.New("entity not found")
	}
	return visitedLocations, nil
}

func (model *VisitedLocation) DeleteVisitedLocationService() error {
	var animalIdCheck int //checking not found err for animal id 404
	err := db.Get(&animalIdCheck, "SELECT id FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	}

	var checkVisitedPointId int // checking if not found location_id in locations 404
	err = db.Get(&checkVisitedPointId, "SELECT id FROM locations WHERE id = $1", model.LocationPointId)
	if err != nil {
		return errors.New("entity not found")
	}

	var checkForExistingTypeIdOnAnimal int // checking if animal not had this typeId 404
	err = db.Get(&checkForExistingTypeIdOnAnimal, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	}

	var getChippedLocationPoint int
	err = db.Get(&getChippedLocationPoint, "SELECT chipping_location_id FROM animals WHERE id = $1", model.Id)

	var idOfFirstVisitedLocation int
	err = db.Get(&idOfFirstVisitedLocation, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1 ORDER BY id LIMIT 1", model.Id)

	if model.LocationPointId == idOfFirstVisitedLocation {
		var firstTwoRows []int // checking if second visitedPoint == chippedPoint
		err = db.Get(&firstTwoRows, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1 ORDER BY id LIMIT 2")
		if err != nil {
			fmt.Println(err.Error())
		} else if firstTwoRows[1] == getChippedLocationPoint {
			_, err = db.Exec("DELETE FROM \"animals-locations\" WHERE animal_id = $1 AND location_id IN ($2, $3)",
				model.Id, firstTwoRows[0], firstTwoRows[1])
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			_, err = db.Exec("DELETE FROM \"animals-locations\" WHERE animal_id = $1 AND location_id = $2",
				model.Id, model.LocationPointId)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	} else {
		_, err = db.Exec("DELETE FROM \"animals-locations\" WHERE animal_id = $1 AND location_id = $2",
			model.Id, model.LocationPointId)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}
