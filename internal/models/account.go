package models

import (
	"errors"
	"fmt"
	"goSimbir/internal/dto"
)

type Account struct {
	Id        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name" json:"lastName"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password,omitempty"`
}

func (model *Account) RegisterAccountService() error {
	err := db.QueryRow("INSERT INTO accounts (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id",
		model.FirstName, model.LastName, model.Email, model.Password).Scan(&model.Id)
	if err != nil {
		return errors.New("account already created")
	}
	return nil
}

func (model *Account) GetAccountService() error {
	err := db.Get(model, "SELECT id, first_name, last_name, email FROM accounts WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("account not found")
	}
	return nil
}

func (model *Account) FindAccountsService(fields dto.AccountFindFields) (*[]Account, error) {
	var accounts []Account
	query := `SELECT id, first_name, last_name, email FROM accounts
         WHERE ($1 = '' OR first_name ILIKE '%'||$1||'%') AND
               ($2 = '' OR last_name ILIKE '%'||$2||'%') AND
               ($3 = '' OR email ILIKE '%'||$3||'%') 
         ORDER BY id LIMIT $4 OFFSET $5`
	err := db.Select(&accounts, query, fields.FirstName, fields.LastName, fields.Email, fields.Size, fields.From)
	if err != nil {
		return nil, err
	}
	return &accounts, nil
}

func (model *Account) UpdateAccountService() error {
	err := db.QueryRow("UPDATE accounts SET first_name = $1, last_name = $2, email = $3 WHERE id = $4 RETURNING id",
		model.FirstName, model.LastName, model.Email, model.Id).Scan(&model.Id)
	if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"accounts_email_key\"" {
		return errors.New("email already exists")
	} else if err != nil {
		return errors.New("account not found")
	}

	return nil
}

func (model *Account) DeleteAccountService() error {
	var accountRelatedWithAnimal int // check if accountId == chipperId
	err := db.Get(&accountRelatedWithAnimal, "SELECT chipper_id FROM animals WHERE chipper_id = $1", model.Id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		return errors.New("invalid value")
	}

	err = db.QueryRow("DELETE FROM accounts WHERE id=$1 RETURNING id", model.Id).Scan(&model.Id)
	if err != nil {
		return errors.New("account not found")
	}

	return nil
}
