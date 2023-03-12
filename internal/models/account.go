package models

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"goSimbir/internal/dto"
)

type Account struct {
	Id        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name" json:"lastName"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password,omitempty"`
}

var db *sqlx.DB

func SetDB(dbConnect *sqlx.DB) { db = dbConnect }

func (model *Account) RegisterAccountService() error {
	var id int
	err := db.QueryRow("INSERT INTO accounts (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id", model.FirstName, model.LastName, model.Email, model.Password).Scan(&id)
	if err != nil {
		return errors.New("account already created")
	} else {
		model.Id = id
	}
	return nil
}

func (model *Account) GetAccountByIdService() error {
	err := db.Get(model, "SELECT * FROM accounts WHERE id = $1", model.Id)
	if err != nil {
		return errors.New("account not found")
	}
	return nil
}

func (model *Account) FindAccountsService(fields dto.FindFields) (*[]Account, error) {
	var accounts []Account
	query := `SELECT id, first_name, last_name, email FROM accounts
         WHERE ($1 = '' OR first_name ILIKE '%'||$1||'%') AND
               ($2 = '' OR last_name ILIKE '%'||$2||'%') AND
               ($3 = '' OR email ILIKE '%'||$3||'%') 
         LIMIT $4 OFFSET $5`
	err := db.Select(&accounts, query, fields.FirstName, fields.LastName, fields.Email, fields.Size, fields.From)
	if err != nil {
		return nil, err
	}
	return &accounts, nil
}

//func (model *Account) GetAccountByIdService(account_id int) (*Account, error) {
//	accounts := &Account{}
//	err := db.Select(&accounts, "SELECT * FROM accounts AS a WHERE a.id = $1", account_id)
//	if err != nil {
//		return nil, err
//	}
//
//	return accounts, nil
//}
//
//func (model *Account) GetAllAccounts() ([]Account, error) {
//	accounts := make([]Account, 0)
//	err := db.Select(&accounts, "SELECT * FROM accounts")
//	if err != nil {
//		return nil, err
//	}
//	return accounts, nil
//}

//

//
//func (model *Account) FindAccount(fields dto.FindFields) ([]Account, error) {
//	accounts := make([]Account, 0)
//
//	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
//	sql, args, err := queryBuilder(builder.Select("*"), fields)
//	if err != nil {
//		return nil, err
//	}
//
//	db.Select(&accounts, sql, args...)
//
//	fmt.Println(accounts)
//	if len(accounts) < 1 {
//		return nil, errors.New("account not find")
//	}
//	return accounts, nil
//}
//
//func queryBuilder(builder sq.SelectBuilder, fields dto.FindFields) (string, []interface{}, error) {
//	cond := make(sq.ILike)
//
//	if fields.FirstName != "" {
//		cond["first_name"] = fields.FirstName
//	}
//
//	if fields.LastName != "" {
//		cond["last_name"] = fields.LastName
//	}
//
//	if fields.Email != "" {
//		cond["email"] = fields.Email
//	}
//
//	if fields.From > 0 {
//		builder.Limit(uint64(fields.From))
//	}
//
//	if fields.Size > 0 {
//		builder.Offset(uint64(fields.Size))
//	}
//
//	return builder.From("accounts").Where(cond).ToSql()
//}
