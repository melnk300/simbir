package models

import (
	"errors"
	"fmt"
	"goSimbir/internal/dto"
	"strings"
	"time"
)

type Animal struct {
	Id                 int        `db:"id" json:"id"`
	AnimalTypes        []int      `db:"animal_types" json:"animalTypes"`
	Weight             float64    `db:"weight" json:"weight"`
	Length             float64    `db:"length" json:"length"`
	Height             float64    `db:"height" json:"height"`
	Gender             string     `db:"gender" json:"gender"`
	LifeStatus         string     `db:"life_status" json:"lifeStatus"`
	ChippingDateTime   string     `db:"chipping_date_time" json:"chippingDateTime"`
	ChipperId          int        `db:"chipper_id" json:"chipperId"`
	ChippingLocationId int        `db:"chipping_location_id" json:"chippingLocationId"`
	VisitedLocations   []int      `db:"visited_locations" json:"visitedLocations"`
	DeathDateTime      *time.Time `db:"death_date_time" json:"deathDateTime"`
}

func (model *Animal) CreateAnimalService() error {
	var chipperId int
	err := db.Get(&chipperId, "SELECT id FROM accounts WHERE id = $1", model.ChipperId)
	if err != nil {
		return errors.New("entity not found")
	}

	var chippingLocationId int
	err = db.Get(&chippingLocationId, "SELECT id FROM locations WHERE id = $1", model.ChippingLocationId)
	if err != nil {
		return errors.New("entity not found")
	}

	var animalTypeIds []int

	subQuery := "("
	for _, value := range model.AnimalTypes {
		subQuery += fmt.Sprintf("%d, ", value)
	}
	subQuery = strings.TrimSuffix(subQuery, ", ")
	subQuery += ")"

	err = db.Select(&animalTypeIds, "SELECT id FROM animal_types WHERE id IN "+subQuery)
	if err != nil {
		return errors.New("entity not found")
	} else if len(animalTypeIds) != len(model.AnimalTypes) {
		return errors.New("entity not found")
	}

	err = db.QueryRow("INSERT INTO animals (weight, length, height, gender, chipper_id, chipping_location_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, chipping_date_time, life_status",
		model.Weight, model.Length, model.Height, model.Gender, model.ChipperId, model.ChippingLocationId).Scan(&model.Id, &model.ChippingDateTime, &model.LifeStatus)
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO \"animals-animal_types\" (animal_id, animal_type_id) VALUES"
	for _, animalType := range model.AnimalTypes {
		query += fmt.Sprintf("(%d, %d), ", model.Id, animalType)
	}

	query = strings.TrimSuffix(query, ", ")

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	return nil
}

func (model *Animal) GetAnimalService() error {
	err := db.Get(model, "SELECT id, weight, length, height, gender, life_status, chipping_date_time, chipper_id, chipping_location_id, death_date_time FROM animals WHERE id=$1", model.Id)
	if err != nil {
		return errors.New("animal not found")
	}

	//err = db.Select(&model.AnimalTypes, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id=$1", model.Id)
	//if err != nil && err.Error() != "sql: no rows in result set" {
	//	panic(err)
	//}
	//
	//err = db.Select(&model.VisitedLocations, "SELECT location_id FROM \"animals-locations\" WHERE animal_id=$1", model.Id)
	//if err != nil && err.Error() != "sql: no rows in result set" {
	//	panic(err)
	//}
	var animalTypeIds []int
	err = db.Select(&animalTypeIds, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id=$1", model.Id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		model.AnimalTypes = nil
	}
	model.AnimalTypes = animalTypeIds

	var visitedLocationIds []int
	err = db.Select(&visitedLocationIds, "SELECT location_id FROM \"animals-locations\" WHERE animal_id=$1", model.Id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		model.VisitedLocations = nil
	}

	model.VisitedLocations = visitedLocationIds

	return nil
}

func (model *Animal) FindAnimalService(fields dto.AnimalFindFields) ([]Animal, error) {
	var animals []Animal
	query := `SELECT id, weight, length, height, gender, life_status, chipping_date_time, chipper_id, chipping_location_id, death_date_time FROM animals
         WHERE ($1 = 0 OR chipper_id = $1) AND
               ($2 = 0 OR chipping_location_id = $2) AND
               ($3 = '' OR life_status ILIKE '%'||$3||'%') AND
               ($4 = '' OR gender ILIKE '%'||$4||'%')`
	if fields.StartDateTime != "" {
		query += fmt.Sprintf(" AND chipping_date_time >= '%s'", fields.StartDateTime)
	}
	if fields.EndDateTime != "" {
		query += fmt.Sprintf(" AND chipping_date_time <= '%s'", fields.EndDateTime)
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", fields.Size, fields.From)

	err := db.Select(&animals, query, fields.ChipperId, fields.ChippingLocationId, fields.LifeStatus, fields.Gender)
	if err != nil {
		return nil, errors.New("animals not found")
	}

	if len(animals) == 0 {
		return nil, errors.New("animals not found")
	}

	for i, animalsStruct := range animals {
		var animalTypeIds []int
		err = db.Select(&animalTypeIds, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id=$1", animalsStruct.Id)
		if err != nil && err.Error() == "sql: no rows in result set" {
			animalsStruct.AnimalTypes = nil
		}
		animals[i].AnimalTypes = animalTypeIds

		var visitedLocationIds []int
		err = db.Select(&visitedLocationIds, "SELECT location_id FROM \"animals-locations\" WHERE animal_id=$1", animalsStruct.Id)
		if err != nil && err.Error() == "sql: no rows in result set" {
			animalsStruct.VisitedLocations = nil
		}
		animals[i].VisitedLocations = visitedLocationIds

	}

	//var animalTypeIds []int
	//err = db.Select(animalTypeIds, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id=$1", model.Id)
	//if err != nil && err.Error() == "sql: no rows in result set" {
	//	model.AnimalTypes = nil
	//}
	//model.AnimalTypes = animalTypeIds
	//
	//var visitedLocationIds []int
	//err = db.Select(&visitedLocationIds, "SELECT location_id FROM \"animals-locations\" WHERE animal_id= $1", model.Id)
	//if err != nil && err.Error() == "sql: no rows in result set" {
	//	model.VisitedLocations = nil
	//}
	//fmt.Println(visitedLocationIds)
	//model.VisitedLocations = visitedLocationIds

	return animals, nil
}

func (model *Animal) UpdateAnimalService() error {
	animal := Animal{}
	err := db.Get(&animal, "SELECT id, life_status, chipping_date_time FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	} else if animal.LifeStatus == "DEAD" && model.LifeStatus == "ALIVE" {
		return errors.New("invalid value")
	}

	err = db.Get(&animal, "SELECT id FROM accounts WHERE id = $1", model.ChipperId)
	if err != nil {
		return errors.New("entity not found")
	}

	err = db.Get(&animal, "SELECT id FROM locations WHERE id = $1", model.ChippingLocationId)
	if err != nil {
		return errors.New("entity not found")
	}

	var firstChippedLocation int
	err = db.Get(&firstChippedLocation, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1 LIMIT 1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	} else if firstChippedLocation == model.ChippingLocationId {
		return errors.New("invalid value")
	}

	query := "UPDATE animals SET weight = $1, length = $2, height = $3, gender = $4, life_status = $5, chipper_id = $6, chipping_location_id = $7"

	if model.LifeStatus == "DEAD" {
		query += ", death_date_time = NOW()"
	}

	query += " WHERE id = $8 RETURNING id, death_date_time, chipping_date_time"

	err = db.QueryRow(query,
		model.Weight, model.Length, model.Height, model.Gender, model.LifeStatus, model.ChipperId, model.ChippingLocationId, model.Id).Scan(&model.Id, &model.DeathDateTime, &model.ChippingDateTime)

	var animalTypeIds []int
	err = db.Select(&animalTypeIds, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id=$1", model.Id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		model.AnimalTypes = nil
	}
	model.AnimalTypes = animalTypeIds

	var visitedLocationIds []int
	err = db.Select(&visitedLocationIds, "SELECT location_id FROM \"animals-locations\" WHERE animal_id=$1", model.Id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		model.VisitedLocations = nil
	}

	model.VisitedLocations = visitedLocationIds

	if err != nil {
		panic(err)
	}
	return nil
}

func (model *Animal) DeleteAnimalService() error {
	var animalId int
	err := db.Get(&animalId, "SELECT id FROM animals WHERE id = $1", model.Id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return errors.New("animal not found")
	}

	var visitedLocationIds []int
	_ = db.Select(&visitedLocationIds, "SELECT location_id FROM \"animals-locations\" WHERE animal_id=$1", model.Id)
	if len(visitedLocationIds) > 0 {
		return errors.New("invalid value")
	}

	_, err = db.Exec("DELETE FROM animals WHERE id=$1", model.Id)
	return nil
}

func (model *Animal) AddAnimalTypeToAnimalService(animalTypeId int) error {
	var animalIdCheck int //checking not found err for animal id
	err := db.Get(&animalIdCheck, "SELECT id FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	}

	var typeIdCheck int //checking not found err for type id 404
	err = db.Get(&typeIdCheck, "SELECT id FROM animal_types WHERE id = $1", animalTypeId)
	if err != nil {
		return errors.New("entity not found")
	}

	//if utils.Contains(animalTypeId, model.AnimalTypes) == true {
	//	return errors.New("typeId already exist for this animalId")
	//}
	//animal.AnimalTypes = append(animal.AnimalTypes, animalTypeId)
	err = db.Get(model, "SELECT weight, length, height, gender, life_status, chipping_date_time, chipper_id, chipping_location_id, death_date_time FROM animals WHERE id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = db.Exec("INSERT INTO \"animals-animal_types\" (animal_id, animal_type_id) VALUES ($1, $2)", model.Id, animalTypeId)
	if err != nil {
		return errors.New("typeId already exist for this animalId")
	}

	err = db.Select(&model.AnimalTypes, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.Select(&model.VisitedLocations, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (model *Animal) UpdateAnimalTypeToAnimalService(editData dto.AnimalEdit) error {
	var animalIdCheck int //checking not found err for animal id
	err := db.Get(&animalIdCheck, "SELECT id FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	}

	var typeIdCheck int //checking not found err for type id 404
	err = db.Get(&typeIdCheck, "SELECT id FROM animal_types WHERE id = $1", editData.OldTypeId)
	if err != nil {
		return errors.New("entity not found")
	}

	err = db.Get(&typeIdCheck, "SELECT id FROM animal_types WHERE id = $1", editData.NewTypeId)
	if err != nil {
		return errors.New("entity not found")
	}

	var animalTypeId int // checking for existing OldTypeId for this animalId OR checking if NewTypeId already exist for this animal (404/409)
	err = db.Get(&animalTypeId, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id = $1 AND animal_type_id IN ($2, $3)",
		model.Id, editData.OldTypeId, editData.NewTypeId)
	//узнать ошибку
	if err != nil {
		return errors.New("entity not found")
	} //else if animalTypeId == editData.NewTypeId {
	//	return errors.New("already exist")
	//}

	_, err = db.Exec("UPDATE \"animals-animal_types\" SET animal_type_id = $1 WHERE animal_type_id = $2",
		editData.NewTypeId, editData.OldTypeId)
	if err != nil {
		return errors.New("already exist")
	}
	err = db.Get(model, "SELECT weight, length, height, gender, life_status, chipping_date_time, chipper_id, chipping_location_id, death_date_time FROM animals WHERE id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.Select(&model.AnimalTypes, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.Select(&model.VisitedLocations, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (model *Animal) DeleteAnimalTypeToAnimal(typeId int) error {
	var animalIdCheck int //checking not found err for animal id 404
	err := db.Get(&animalIdCheck, "SELECT id FROM animals WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	}

	var typeIdCheck int //checking not found err for type id 404
	err = db.Get(&typeIdCheck, "SELECT id FROM animal_types WHERE id = $1", typeId)
	if err != nil {
		return errors.New("entity not found")
	}

	var checkForExistingTypeIdOnAnimal int // checking if animal not had this typeId 404
	err = db.Get(&checkForExistingTypeIdOnAnimal, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id = $1", model.Id)
	if err != nil {
		return errors.New("entity not found")
	}

	var checkIfOneRowAndItTypeId []int // checking for one row in types of animal id and if its == typeId then 400
	err = db.Get(checkIfOneRowAndItTypeId, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	} else if len(checkIfOneRowAndItTypeId) == 1 && checkIfOneRowAndItTypeId[0] == typeId {
		return errors.New("invalid value")
	}

	_, err = db.Exec("DELETE FROM \"animals-animal_types\" WHERE animal_id = $1 AND animal_type_id = $2", model.Id, typeId)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.Get(model, "SELECT weight, length, height, gender, life_status, chipping_date_time, chipper_id, chipping_location_id, death_date_time FROM animals WHERE id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.Select(&model.AnimalTypes, "SELECT animal_type_id FROM \"animals-animal_types\" WHERE animal_id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.Select(&model.VisitedLocations, "SELECT location_id FROM \"animals-locations\" WHERE animal_id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
